package cp_app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"runtime/debug"
	"time"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_middleware"
	"warehouse/v5-go-component/cp_obj"
	"warehouse/v5-go-component/cp_tracing"
)

type Http struct {
	svr		*Ins
	BaseServer	*http.Server //gin本体
	handler		*Handler
}

// InitRouter 初始化路由
func (this *Http) InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.MaxMultipartMemory = 100 * 1024

	//router.Use(gin.LoggerWithWriter(os.Stdout))
	router.Use(cp_middleware.Recovery(this.svr.Limiter))
	router.Use(cp_middleware.InLimiter(this.svr.Limiter))

	router.GET("/heartbeat", this.handler.HeartBeat())

	//g1 := router.Group("/faciliodata", gin.BasicAuth(gin.Accounts{
	//	"foo":    "bar",
	//	"austin": "1234",
	//	"lena":   "hello2",
	//	"manu":   "4321",
	//}))
	//{
	//	g1.POST("", this.handler.Faciliodata())
	//}

	rv1 := router.Group("/api/v2", cp_middleware.Prepare())
	{
		rv2 := rv1.Group("/info")
		{
			rv2.GET("/:svrName", this.handler.DispatchInfo())
		}

		rv3 := rv1.Group("/common")
		{
			rv3.Any("/:svrName/:controller/:funName", this.handler.DispatchApi(false))
		}

		rv4 := rv1.Group("/admin")
		{
			rv4.Any("/:svrName/:controller/:funName", this.handler.DispatchApi(true))
		}

		rv5 := rv1.Group("/special")
		{
			rv5.Any("/:svrName/:controller/:funName/:" + cp_constant.SPECIAL_ID, this.handler.DispatchApi(false))
		}
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(200, cp_obj.NewResponse().Err("no router found!"))
	})

	return router
}

func (this *Http) StartSvr() {
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
				err := Instance.TraceLog.PushRuntime(cp_tracing.NewTraceRuntime(this.svr.ServerName, cp_constant.TracingLevelCritical, msg))
				if err != nil {
					cp_log.Error("RuntimePush:" + err.Error())
				}
			}

			os.Exit(-120)
		}
	}()

	err := this.BaseServer.ListenAndServe()
	if err != nil {
		cp_log.Error("服务监听出错:" + err.Error())
	}

	this.svr.stopChan <- err
}

func NewHttp(svr *Ins) *Http {
	h := &Http{
		svr:		svr,
		handler: 	NewHandler(svr),
		BaseServer:	&http.Server {
			Addr:		fmt.Sprintf(":%d", svr.DataCenter.Base.HttpPort),
			ReadTimeout:    60 * time.Second,
			WriteTimeout:   60 * time.Second,
		},
	}
	h.BaseServer.Handler = h.InitRouter() //return gin router

	return h
}


