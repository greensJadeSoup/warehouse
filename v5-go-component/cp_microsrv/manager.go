package cp_microsrv

import (
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"warehouse/v5-go-component/cp_dc"
	"github.com/hashicorp/consul/api"
)

//serviceHeartbeatPath HTTP服务心跳检查路径、TCP服务连接检查标志
const HEARTBEAT_PATH = "/heartbeat"

type MicroSvrManager struct {
	isLocal		bool
	Registed	bool

	consulClient       *api.Client
	consulQueryOptions *api.QueryOptions

	mu           sync.Mutex
	table        atomic.Value
}


func (m *MicroSvrManager) getTable() map[string]SvrTable {
	return m.table.Load().(map[string]SvrTable)
}

func (m *MicroSvrManager) setTable(st map[string]SvrTable) {
	if st == nil {
		return
	}
	m.mu.Lock()
	m.table.Store(st)
	m.mu.Unlock()
}

//Register 注册服务
func (m *MicroSvrManager) Register(svrName, svrID string, appConf cp_dc.DcBaseConfig) (err error) {
	agent := m.consulClient.Agent()

	//ServiceID Format：type-name-address-httPort-rpcPort
	svcReg := &api.AgentServiceRegistration{
		ID:      svrID,
		Name:    svrName,
		Address: appConf.IP,
		Port:    appConf.HttpPort,
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d%s", appConf.IP, appConf.HttpPort, HEARTBEAT_PATH),
			Interval: "1s",
			Timeout: "3s",
		},
	}

	if err = agent.ServiceRegister(svcReg); err != nil {
		err = fmt.Errorf("注册服务[%s]失败：%s", svrName, err)
		return
	}

	cp_log.Info("[MicroSvrManager]server register success, checkID:" + svrID)
	m.Registed = true
	return
}

//UnRegister 反注册服务
func (m *MicroSvrManager) UnRegister(svrID string) error {
	if !m.Registed {
		return nil
	}

	agent := m.consulClient.Agent()

	if err := agent.ServiceDeregister(svrID); err != nil {
		return fmt.Errorf("consul：反注册服务失败：%s", err)
	}
	cp_log.Warning("Consul注销成功, checkID: " + svrID)

	return nil
}

//Watch 监控服务健康，后台执行
func (m *MicroSvrManager) Watch() {
	defer func() {
		if err := recover(); err != nil {
			cp_log.Error(fmt.Sprintf("stack: %s", debug.Stack()))
		}
	}()

	for {
		err := m.checkHealth()
		if err != nil {
			cp_log.Error(err.Error())
		}
		time.Sleep(10 * time.Second) //休眠片刻
	}
}

//检查健康服务
//ServiceID Format：type-name-address-httPort-rpcPort
func (m *MicroSvrManager) checkHealth() (error) {
	checks, _, err := m.consulClient.Health().State(api.HealthPassing, m.consulQueryOptions)
	if err != nil {
		return cp_error.NewSysError(fmt.Errorf("检查健康服务出错：%s", err))
	}

	tableMap := make(map[string]SvrTable)

	for i, l := 0, len(checks); i < l; i++ {
		sid := strings.ToLower(checks[i].ServiceID)
		if sid == "" {
			continue
		}

		svrInstance := MicroService{MicroSvrMgr: m}
		err = svrInstance.Init(sid)
		if err != nil {
			continue
		}
		st, ok := tableMap[svrInstance.Name]
		if !ok {
			st = SvrTable{}
		}

		st.Array = append(st.Array, svrInstance)
		tableMap[svrInstance.Name] = st
	}

	m.setTable(tableMap)
	return nil
}

//GetServiceList 获取服务列表
func (m *MicroSvrManager) GetServiceToCall(name string) (*MicroService) {
	st := m.getTable()
	stb, ok := st[strings.ToLower(name)]
	if !ok {
		return nil
	}

	sl := len(stb.Array)
	if sl == 0 {
		return nil
	}

	if sl > 1 && stb.Index == sl {
		stb.Index = 0
	}

	CallIndex := stb.Index
	stb.Index ++

	return &(stb.Array[CallIndex])
}

// v5-go-component
// 获取第三方包构造的client
func (m *MicroSvrManager) ConsulClient() *api.Client {
	return m.consulClient
}

// v5-go-component
// 创建Consul客户端
func NewConsulClient(DcAddress string, token string) (*api.Client, error) {
	if DcAddress == "" {
		return nil, cp_error.NewSysError("API Service Manager New Failed: DcAddress Invalid")
	}

	consulConf := api.DefaultConfig()
	consulConf.Address = DcAddress
	consulConf.Token = token

	consulClient, err := api.NewClient(consulConf)
	if err != nil {
		return nil, cp_error.NewSysError("API Service Manager New Failed: " + err.Error())
	}

	return consulClient, nil
}

func NewMicroSvrManagerByConsulClient(consulClient *api.Client, isLocal bool) (*MicroSvrManager, error) {
	if consulClient == nil {
		return nil, cp_error.NewSysError("API Service Manager New Failed: ConsulClient Invalid")
	}

	m := &MicroSvrManager{
		isLocal: 	    isLocal,
		consulClient:       consulClient,
		consulQueryOptions: &api.QueryOptions{RequireConsistent: true, WaitIndex: 0},
	}

	err := m.checkHealth()
	if err != nil {
		return nil, err
	}

	return m, nil
}