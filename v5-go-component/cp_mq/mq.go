package cp_mq

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	mqInstance = &mq{}
)

type mq struct {
	consumerAdapter map[string]reflect.Type
	producerAdapter map[string]reflect.Type
	consumerList    []IConsumer
	producerList    []IProducer
}

//注册消费者适配器
func ConsumerRegister(adapter string, ci IConsumer) {
	cv := reflect.ValueOf(ci)
	ct := reflect.Indirect(cv).Type()

	if mqInstance.consumerAdapter == nil {
		mqInstance.consumerAdapter = make(map[string]reflect.Type)
	}
	mqInstance.consumerAdapter[strings.ToLower(adapter)] = ct
}

//注册发布者适配器
func ProducerRegister(adapter string, pi IProducer) {
	pv := reflect.ValueOf(pi)
	pt := reflect.Indirect(pv).Type()

	if mqInstance.producerAdapter == nil {
		mqInstance.producerAdapter = make(map[string]reflect.Type)
	}

	mqInstance.producerAdapter[strings.ToLower(adapter)] = pt
}

func Exit() error {
	var (
		err    error
		errMsg string
	)

	for _, p := range mqInstance.producerList {
		err = p.Stop()
		if err != nil {
			if len(errMsg) > 0 {
				errMsg += ","
			}
			errMsg += fmt.Sprintf("%q", err)
		}
	}

	for _, c := range mqInstance.consumerList {
		err = c.Stop(errors.New("进程主动退出"))
		if err != nil {
			if len(errMsg) > 0 {
				errMsg += ","
			}
			errMsg += fmt.Sprintf("%q", err)
		}
	}

	if len(errMsg) > 0 {
		err = fmt.Errorf("MQ: Exit error %s", errMsg)
	}

	return err
}
