package consumer

import (
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_log"
)

// 订阅者在会话中消费消息，并标记当前消息已经被消费。
func KAFKAMessageHandler1(message string) (error, int) {
	//cp_log.Info().Msgf("[Kafka1]消息消费成功: value = %s", message)
	return nil, 0
}

// 订阅者在会话中消费消息，并标记当前消息已经被消费。
func RocketMQMessageHandler1(message string) (error, cp_constant.MQ_ERR_TYPE) {
	cp_log.Info("enter")
	cp_log.Info("[RocketMQ]消息消费成功:" + message)
	return nil, cp_constant.MQ_ERR_TYPE_OK
}





