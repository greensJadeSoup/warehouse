package cp_app

import (
	"fmt"
	"golang.org/x/net/context"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"
	"warehouse/v5-go-component/cp_cache"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_dc"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_ini"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_middleware"
	"warehouse/v5-go-component/cp_mq"
	_ "warehouse/v5-go-component/cp_mq/kafka"
	_ "warehouse/v5-go-component/cp_mq/nsq"
	_ "warehouse/v5-go-component/cp_mq/rmq"
	"warehouse/v5-go-component/cp_nosql"
	"warehouse/v5-go-component/cp_orm"
	"warehouse/v5-go-component/cp_task"
	"warehouse/v5-go-component/cp_tracing"
	"warehouse/v5-go-component/cp_util"
)

var Instance *Ins

type Ins struct {
	SvrVersion	string
	CpVersion	string
	ServerName	string
	ID		string

	stopChan	chan error
	signalChan	chan os.Signal

	LocalConf	*cp_ini.LocalConfig

	Logger		*cp_log.Logger
	TaskManager	*cp_task.TaskManager

	Mongo		*cp_nosql.Mongo

	HttpSvr		*Http

	DataCenter	*cp_dc.DcConfig
	TraceLog	*cp_tracing.Tracer

	CallClient	*CallClient

	Limiter		*cp_middleware.Limiter
	Info		*Info
}

//initInstance 初始化APIServer实例
func InitInstance(svrName, svrVersion string, regInitFun... func()) {
	defer func() {
		if err := recover(); err != nil {
			var msg string

			switch e := err.(type) {
			case error:
				msg = "[Error]: " + e.Error() + " [Stack]:" + string(debug.Stack())
			case string:
				msg = "[Error]: " + e + " [Stack]:" + string(debug.Stack())
			}
			cp_log.Error(msg)

			if Instance.TraceLog != nil {
				err := Instance.TraceLog.PushRuntime(cp_tracing.NewTraceRuntime(svrName, cp_constant.TracingLevelCritical, msg))
				if err != nil {
					cp_log.Error("RuntimePush:" + err.Error())
				}
			}

			os.Exit(-100)
		}
	}()

	conf, err := cp_ini.GetLocalConfig()
	if err != nil {
		panic("GetLocal 失败：" + err.Error())
	}

	err = NewInstance(svrName, svrVersion, conf)
	if err != nil {
		panic("Instance 创建失败：" + err.Error())
	}

	for _, f := range regInitFun {
		f()
	}
}

func NewInstance(svrName, svrVersion string, conf *cp_ini.LocalConfig) error {
	var err error

	svr := &Ins{
		LocalConf:		conf,
		stopChan:		make(chan error, 1),
		signalChan:		make(chan os.Signal, 1),
		CpVersion:		cp_constant.ComponentVersion,
	}
	Instance = svr

	signal.Notify(svr.signalChan, syscall.SIGTERM)
	signal.Notify(svr.signalChan, syscall.SIGKILL)
	signal.Notify(svr.signalChan, syscall.SIGQUIT)
	signal.Notify(svr.signalChan, syscall.SIGINT)
	signal.Notify(svr.signalChan, syscall.SIGHUP)

	//根据consulClient创建Consul数据配置中心的配置对象
	svr.DataCenter, err = cp_dc.NewDataCenter(conf)
	if err != nil {
		return err
	}

	//初始化Snowflake算法
	cp_util.InitSnowflake(1,1)

	//开启日志，顺序需在靠前，不可调整
	svr.Logger = cp_log.NewLogger(&svr.DataCenter.Log)
	svr.SetService(svrName, svrVersion)

	svr.Limiter = cp_middleware.NewLimiter(svr.DataCenter.Base.ServingLimit)
	//svr.TraceLog = cp_tracing.NewTracer(svr.DataCenter)
	svr.TaskManager = cp_task.GetTaskManager()
	svr.CallClient = NewCallClient(svr)
	svr.HttpSvr = NewHttp(svr)
	svr.Info = NewInfo(svr)

	if _, err = cp_orm.NewEngine(&svr.DataCenter.DBList); err != nil {
		return err
	}

	if err = cp_cache.InitCache(&svr.DataCenter.Cache); err != nil {
		return err
	}

	return nil
}

func (svr * Ins) Start() error {
	//启动多例任务
	err := svr.TaskManager.StartMultiple()
	if err != nil {
		return cp_error.NewSysError("多例任务启动失败：" + err.Error())
	}
	cp_log.Info("多例任务列表开始执行...")

	//启动单例任务
	if svr.DataCenter.Base.IsLeader == true {
		err = svr.TaskManager.StartSingle()
		if err != nil {
			return cp_error.NewSysError("单例任务启动失败：" + err.Error())
		}
		cp_log.Info("单例任务列表开始执行...")
	}

	//go svr.MicroSvrManager.Watch()
	go svr.HttpSvr.StartSvr()

	cp_log.Info(fmt.Sprintf("服务[%s]已开始监听http端口: %d", svr.ServerName, svr.DataCenter.Base.HttpPort))

	select {
	case s := <- svr.signalChan:
		err = cp_error.NewSysError(fmt.Sprintf("[%s-%s] 捕获到退出信号: %v", svr.ServerName, svr.SvrVersion, s.String()))
	case err = <- svr.stopChan:
		err = cp_error.NewSysError(fmt.Sprintf("[%s-%s] 服务监听出错: %v", svr.ServerName, svr.SvrVersion, err))
	}

	return err
}

func (svr * Ins) Stop() {
	var wg sync.WaitGroup

	//关闭http服务
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()
	err := svr.HttpSvr.BaseServer.Shutdown(ctx)
	if err != nil {
		cp_log.Error("关闭http链路追踪失败:" + err.Error())
	}
	<-ctx.Done()

	//关闭mongodb
	if svr.Mongo != nil {
		wg.Add(1)
		go func() {
			err := svr.Mongo.Client.Disconnect(svr.Mongo.Ctx)
			if err != nil {
				cp_log.Error("MongoDB服务退出失败："+err.Error())
			} else {
				cp_log.Info("MongoDB服务退出成功")
			}
			wg.Done()
		}()
	}

	//关闭任务列表
	wg.Add(1)
	go func() {
		cp_log.Info("停止任务中心")
		svr.TaskManager.Stop()
		svr.TaskManager.Wait()
		cp_log.Info("任务中心已停止")
		wg.Done()
	}()

	//关闭MQ
	wg.Add(1)
	go func() {
		cp_log.Info("停止所有MQ")
		err := cp_mq.Exit()
		if err != nil {
			cp_log.Error("MQ服务退出失败："+err.Error())
		} else {
			cp_log.Info("MQ服务退出成功")
		}

		wg.Done()
	}()

	//等待
	wg.Wait()

	cp_log.Info(fmt.Sprintf(fmt.Sprintf("[%s-%s-%s] 服务已完全退出", svr.ServerName, svr.SvrVersion, svr.CpVersion)))
}

func (this * Ins) SetService(serviceName, serviceVersion string) {
	this.SvrVersion = serviceVersion
	this.ServerName = serviceName

	this.ID = fmt.Sprintf("%s-%s-%s-%d",
		"go",
		this.ServerName,
		this.DataCenter.Base.IP,
		this.DataCenter.Base.HttpPort)
	cp_log.Info(fmt.Sprintf("服务版本[%s] 组件版本[%s]", this.SvrVersion, this.CpVersion))
}

//func (this * Ins) SetErrEnv(flag bool) {
//	cp_error.SetErrEnv(flag)
//}

func GetIns() *Ins {
	return Instance
}

//Run 创建并启动 Instance 服务
func Run() {
	defer func() {
		if err := recover(); err != nil {
			var msg string

			switch e := err.(type) {
			case error:
				msg = "[Error]: " + e.Error() + " [Stack]:" + string(debug.Stack())
			case string:
				msg = "[Error]: " + e + " [Stack]:" + string(debug.Stack())
			}
			cp_log.Error(msg)

			if Instance.TraceLog != nil {
				err := Instance.TraceLog.PushRuntime(cp_tracing.NewTraceRuntime(Instance.ServerName, cp_constant.TracingLevelCritical, msg))
				if err != nil {
					cp_log.Error("RuntimePush:" + err.Error())
				}
			}

			os.Exit(-110)
		}
	}()

	err := Instance.Start()
	if err != nil {
		Instance.Stop()
		panic("服务异常退出:" + err.Error())
	}
}
