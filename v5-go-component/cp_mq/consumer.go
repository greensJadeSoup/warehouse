package cp_mq

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"warehouse/v5-go-component/cp_constant"
)

const (
	CONSUMER_STATE_INIT ConsumerState = iota //初始化状态
	CONSUMER_STATE_STOP 			 //停止接收消息
	CONSUMER_STATE_STOPING  		 //停止中
	CONSUMER_STATE_RUN                       //消息处理中
	CONSUMER_STATE_WAIT                      //等待消息投递

)

type ConsumerState int

func (c ConsumerState) String() string {
	return consumerStateToString(c)
}

func consumerStateToString(s ConsumerState) string {
	str := ""
	switch s {
	case CONSUMER_STATE_WAIT:
		str = "等待"
	case CONSUMER_STATE_RUN:
		str = "执行"
	case CONSUMER_STATE_STOP:
		str = "停止"
	case CONSUMER_STATE_STOPING:
		str = "正在停止中"
	default:
		str = fmt.Sprintf("未知[%v]", s)
	}

	return str
}

type ConsumerHandlerFunc func(string) (error, cp_constant.MQ_ERR_TYPE)

type IConsumer interface {
	Init(configStr string, handler ConsumerHandlerFunc) error
	Start() error
	Stop(...error) error
	State() ConsumerState
	SetComsumerStatus(s ConsumerState)
	Error() error
	Ctx() context.Context

	Name() string
	RunCount() int64
	DelayAvg() int64
	AddRunCount()
	AddDelaySum(d int64)
}

func NewConsumer(adapter, configStr string, handler ConsumerHandlerFunc) (IConsumer, error) {
	if handler == nil {
		return nil, errors.New("MQ：handler func cannot be nil")
	}

	v, ok := mqInstance.consumerAdapter[strings.ToLower(adapter)]
	if !ok {
		return nil, fmt.Errorf("MQ: unknown Consumer adapter %q", adapter)
	}

	vo := reflect.New(v)
	ci, ok := vo.Interface().(IConsumer)
	if !ok {
		return nil, fmt.Errorf("MQ: %q is a invalid IConsumer", adapter)
	}

	err := ci.Init(configStr, handler)
	if err != nil {
		return nil, err
	}

	err = ci.Start()
	if err != nil {
		return nil, err
	}

	mqInstance.consumerList = append(mqInstance.consumerList, ci)

	return ci, nil
}

//获得消费者列表
func GetConsumerList() []IConsumer {
	return mqInstance.consumerList
}
