package mq

import (
	"fmt"
	"strings"
	"time"
	"warehouse/v5-go-api-shopee/bll"
	"warehouse/v5-go-api-shopee/conf"
	"warehouse/v5-go-api-shopee/mq/consumer"
	"warehouse/v5-go-api-shopee/mq/producer"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_mq"
)

func InitMQ(){
	initKafka()
	//initRocketMQ()
}

//开始创建kafka订阅者
func initKafka() {
	kafkaConf, err := cp_app.GetIns().DataCenter.GetMQ("kafka")
	if err != nil {
		panic(err)
	}

	kafkaAddr := strings.Join(kafkaConf.Server, `","`)
	cp_log.Info(kafkaAddr)

	producer.ProducerSyncOrderTask, err = cp_mq.NewProducer(
		"kafka",
		fmt.Sprintf(`{"address":["%s"],"topic":"v2_Sync_Order"}`, kafkaAddr))
	if err != nil {
		panic(err)
	}

	producer.ProducerSyncItemAndModelTask, err = cp_mq.NewProducer(
		"kafka",
		fmt.Sprintf(`{"address":["%s"],"topic":"v2_Sync_Item_And_Model"}`, kafkaAddr))
	if err != nil {
		panic(err)
	}

	producer.ProducerSyncPushOrderStatusTask, err = cp_mq.NewProducer(
		"kafka",
		fmt.Sprintf(`{"address":["%s"],"topic":"v2_Sync_Push_Order_Status"}`, kafkaAddr))
	if err != nil {
		panic(err)
	}

	if !conf.GetAppConfig().MQConsumerSyncOrder.Stop {
		_, err = cp_mq.NewConsumer(
			"kafka",
			fmt.Sprintf(`{"address":["%s"],"topic":"v2_Sync_Order","group":"v2_Sync_Order_Group_1","need_dlq":true}`, kafkaAddr),
			bll.NewOrderBL(nil).ConsumerOrder)
		if err != nil {
			panic(err)
		}
	}

	if !conf.GetAppConfig().MQConsumerSyncItem.Stop {
		_, err = cp_mq.NewConsumer(
			"kafka",
			fmt.Sprintf(`{"address":["%s"],"topic":"v2_Sync_Item_And_Model","group":"v2_Sync_Item_And_Model_Group_1","need_dlq":true}`, kafkaAddr),
			bll.NewItemBL(nil).ConsumerItemAndModel)
		if err != nil {
			panic(err)
		}
	}

	if !conf.GetAppConfig().MQConsumerSyncPushOrderStatus.Stop {
		_, err = cp_mq.NewConsumer(
			"kafka",
			fmt.Sprintf(`{"address":["%s"],"topic":"v2_Sync_Push_Order_Status","group":"v2_Sync_Push_Order_Status_Group_1","need_dlq":true}`, kafkaAddr),
			bll.NewOrderBL(nil).ConsumerOrderStatus)
		if err != nil {
			panic(err)
		}
	}
}

func initRocketMQ() {
	var err error

	rocketMQConf, err := cp_app.GetIns().DataCenter.GetMQ("rocketMQ")
	if err != nil {
		panic(err)
	}

	rocketMQAddr := strings.Join(rocketMQConf.Server, `","`)
	cp_log.Info(rocketMQAddr)

	producer.ProducerRocketMQ, err = cp_mq.NewProducer(
		"rocketMQ",
		fmt.Sprintf(`{"address":["%s"],"topic":"testTopic01"}`, rocketMQAddr))
	if err != nil {
		panic(err)
	}

	_, err = cp_mq.NewConsumer(
		"rocketMQ",
		fmt.Sprintf(`{"address":["%s"],"topic":"testTopic01","group":"MyConsumerGroupName"}`, rocketMQAddr),
		consumer.RocketMQMessageHandler1)
	if err != nil {
		panic(err)
	}
}


func ProduceRocketMQ() {
	for {
		time.Sleep(1*time.Second)
		fmt.Println(producer.ProducerRocketMQ.Publish([]byte("2323"), "111"))
	}
}


