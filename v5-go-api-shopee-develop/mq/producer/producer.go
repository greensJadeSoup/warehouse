package producer

import (
	"warehouse/v5-go-component/cp_mq"
)

var ProducerSyncOrderTask cp_mq.IProducer
var ProducerSyncPushOrderStatusTask cp_mq.IProducer
var ProducerSyncItemAndModelTask cp_mq.IProducer
var ProducerRocketMQ cp_mq.IProducer