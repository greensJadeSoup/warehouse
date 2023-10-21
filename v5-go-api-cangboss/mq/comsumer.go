package mq

import (
	"time"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_log"
)

// 订阅者在会话中消费消息，并标记当前消息已经被消费。
func NSQMessageHandler1(message string) (error, cp_constant.MQ_ERR_TYPE) {
	cp_log.Info("[NSQ]消息消费成功:" + message)
	return nil, cp_constant.MQ_ERR_TYPE_OK
}

// 订阅者在会话中消费消息，并标记当前消息已经被消费。
func KAFKAMessageHandler1(message string) (error, cp_constant.MQ_ERR_TYPE) {
	cp_log.Info("[Kafka]消息消费成功:" + message)
	return nil, cp_constant.MQ_ERR_TYPE_OK
}

// 订阅者在会话中消费消息，并标记当前消息已经被消费。
func RocketMQMessageHandler1(message string) (error, cp_constant.MQ_ERR_TYPE) {
	cp_log.Info("enter")
	time.Sleep(1*time.Second)
	cp_log.Info("[RocketMQ]消息消费成功:" + message)
	return nil, cp_constant.MQ_ERR_TYPE_OK
}






