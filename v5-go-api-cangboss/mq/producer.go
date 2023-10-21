package mq

import (
	"fmt"
	"strconv"
	"time"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_mq"
)

var NSQProducerT1 cp_mq.IProducer
var KafkaProducerT1 cp_mq.IProducer
var RocketMQProducerT1 cp_mq.IProducer

func NSQProducter() error {
	for i:=0; i<1000; i++ {
		msg := "hello world" + strconv.Itoa(i)
		err := NSQProducerT1.Publish([]byte(msg), "")
		if err != nil {
			cp_log.Error("send message err=" + err.Error())
			return err
		}

		cp_log.Info("[NSQ]消息发送成功:" + msg)
		time.Sleep(1*time.Second)
	}

	return nil
}

func KafkaProducer() error {
	for i:=0; i<1000; i++ {
		if KafkaProducerT1.State() == cp_mq.PRODUCER_STATE_STOPING ||
			KafkaProducerT1.State() == cp_mq.PRODUCER_STATE_STOP { //收到退出信号，主动退出
			cp_log.Warning("producer worker协程收到退出信号, 主动退出")
			return nil
		}
		msg := "hello world" + strconv.Itoa(i)
		err := KafkaProducerT1.Publish([]byte(msg), "")
		if err != nil {
			cp_log.Error("send message err=" + err.Error())
			return err
		}

		cp_log.Info("[Kafka]消息发送成功:" + msg)
		time.Sleep(1*time.Second)
	}
	fmt.Println("KafkaProducter exit!")
	return nil
}


func RocketMQProducer() error {
	for i:=0; i<1000; i++ {
		if RocketMQProducerT1.State() == cp_mq.PRODUCER_STATE_STOPING || RocketMQProducerT1.State() == cp_mq.PRODUCER_STATE_STOP {
			//收到退出信号，主动退出
			cp_log.Warning("producer worker协程收到退出信号, 主动退出")
			return nil
		}
		msg := "hello world" + strconv.Itoa(i)
		err := RocketMQProducerT1.Publish([]byte(msg), "test01")
		if err != nil {
			cp_log.Error("send message err=" + err.Error())
			return err
		}

		cp_log.Info("[RocketMQ]消息发送成功:" + msg)
		time.Sleep(300*time.Millisecond)
	}
	fmt.Println("KafkaProducter exit!")
	return nil
}
