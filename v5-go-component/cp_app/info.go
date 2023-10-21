package cp_app

import (
	"warehouse/v5-go-component/cp_mq"
	"warehouse/v5-go-component/cp_obj"
)

//Info 服务的运行信息
type Info struct {
	Svr		*Ins 	`json:"-"`   //服务信息

	RunState	map[string]interface{} //运行信息
	TaskState	map[string]interface{} //任务信息
	MqInfo		infoMQ			//消息队列信息
}

//NewInfo 创建新的运行信息对象
func NewInfo(svr *Ins) *Info {
	info := &Info{
		Svr: 		svr,
		RunState:  make(map[string]interface{}),
	}

	info.RunState["SvrVersion"] = svr.SvrVersion
	info.RunState["CpVersion"] = svr.CpVersion
	info.RunState["SvrName"] = svr.ServerName
	info.RunState["StartTime"] = cp_obj.NewDatetime()

	return info
}

type infoMQ struct {
	Version  string
	Producer struct {
		Start int
		Items []infoMQProducerItem
	}
	Consumer struct {
		Start int
		Items []infoMQConsumerItem
	}
}

type infoMQConsumerItem struct {
	Name     string
	Status   string
	RunCount int64
	DelayAvg int64
	Error    string
}

type infoMQProducerItem struct {
	Name     string
	Status   string
	RunCount int64
	DelayAvg int64
}

func (i *Info) setRunState() {
	i.RunState["ServeingLimit"] = i.Svr.DataCenter.Base.ServingLimit
	i.RunState["Serving"] = i.Svr.Limiter.Count.Serving
	i.RunState["Served"] = i.Svr.Limiter.Count.Served
	i.RunState["Reject"] = i.Svr.Limiter.Count.Reject
	i.RunState["Panic"] = i.Svr.Limiter.Count.Panic
	i.RunState["ServedMax"] = i.Svr.Limiter.Count.ServedMax
}

func (i *Info) setComponentMQ() {
	mqInfo := &i.MqInfo

	producerList := cp_mq.GetProducerList()
	producerListLen := len(producerList)

	mqInfo.Producer.Start = producerListLen
	mqInfo.Producer.Items = make([]infoMQProducerItem, producerListLen)

	for i := 0; i < producerListLen; i++ {
		mqInfo.Producer.Items[i].Name = producerList[i].Name()
		mqInfo.Producer.Items[i].Status = producerList[i].State().String()
		mqInfo.Producer.Items[i].RunCount = producerList[i].RunCount()
		mqInfo.Producer.Items[i].DelayAvg = producerList[i].DelayAvg()
	}

	consumerList := cp_mq.GetConsumerList()
	consumerListLen := len(consumerList)

	mqInfo.Consumer.Start = consumerListLen
	mqInfo.Consumer.Items = make([]infoMQConsumerItem, consumerListLen)

	for i := 0; i < consumerListLen; i++ {
		errMsg := "nil"
		if consumerList[i].Error() != nil {
			errMsg = consumerList[i].Error().Error()
		}

		mqInfo.Consumer.Items[i].Name = consumerList[i].Name()
		mqInfo.Consumer.Items[i].Status = consumerList[i].State().String()
		mqInfo.Consumer.Items[i].RunCount = consumerList[i].RunCount()
		mqInfo.Consumer.Items[i].DelayAvg = consumerList[i].DelayAvg()
		mqInfo.Consumer.Items[i].Error = errMsg
	}
}

func (i *Info) setComponentTask() {
	i.TaskState = i.Svr.TaskManager.Info()
}

//JSON 获取info信息JSON化内容
func (i *Info) JSON() *Info {
	var newInfo Info //防止并发操作

	newInfo.Svr = i.Svr
	newInfo.RunState =  make(map[string]interface{})

	newInfo.setRunState()
	newInfo.setComponentMQ()
	newInfo.setComponentTask()

	return &newInfo
}
