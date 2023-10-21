package mq

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"os"
	"strings"
	"time"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_mq"
)

//开始创建kafka订阅者
func InitMQ() {
	//initNSQ()
	initKafka()
	//initRocketMQ()
}

func initNSQ() {
	nsqConf, err := cp_app.GetIns().DataCenter.GetMQ("nsq")
	if err != nil {
		panic(err)
	}

	nsqAddr := strings.Join(nsqConf.Server, `","`)
	cp_log.Info(nsqAddr)

	NSQProducerT1, err = cp_mq.NewProducer(
		"nsq",
		fmt.Sprintf(`{"address":"%s","topic":"topicName05"}`, nsqAddr))
	if err != nil {
		panic(err)
	}

	_, err = cp_mq.NewConsumer(
		"nsq",
		fmt.Sprintf(`{"address":"%s","topic":"topicName05","channel":"ch1"}`, nsqConf.Server[0]),
		NSQMessageHandler1)
	if err != nil {
		panic(err)
	}
}

func initKafka() {
	kafkaConf, err := cp_app.GetIns().DataCenter.GetMQ("kafka")
	if err != nil {
		panic(err)
	}

	kafkaAddr := strings.Join(kafkaConf.Server, `","`)
	cp_log.Info(kafkaAddr)

	KafkaProducerT1, err = cp_mq.NewProducer(
		"kafka",
		fmt.Sprintf(`{"address":["%s"],"topic":"topicName"}`, kafkaAddr))
	if err != nil {
		panic(err)
	}

	_, err = cp_mq.NewConsumer(
		"kafka",
		fmt.Sprintf(`{"address":["%s"],"topic":"topicName","group":"hh-g3"}`, kafkaAddr),
		KAFKAMessageHandler1)
	if err != nil {
		panic(err)
	}
}

func initRocketMQ() {
	rocketMQConf, err := cp_app.GetIns().DataCenter.GetMQ("rocketMQ")
	if err != nil {
		panic(err)
	}

	rocketMQAddr := strings.Join(rocketMQConf.Server, `","`)
	cp_log.Info(rocketMQAddr)

	RocketMQProducerT1, err = cp_mq.NewProducer(
		"rocketMQ",
		fmt.Sprintf(`{"address":["%s"],"topic":"testTopic04"}`, rocketMQAddr))
	if err != nil {
		panic(err)
	}

	_, err = cp_mq.NewConsumer(
		"rocketMQ",
		fmt.Sprintf(`{"address":["%s"],"topic":"testTopic04","group":"MyConsumerGroupName"}`, rocketMQAddr),
		RocketMQMessageHandler1)
	if err != nil {
		panic(err)
	}
}

func SendSyncMessage(message string) {
	// 发送消息
	endPoint := []string{"192.168.10.17:9876"}
	// 创建一个producer实例
	p, _ := rocketmq.NewProducer(
		producer.WithNameServer(endPoint),
		producer.WithRetry(2),
		producer.WithGroupName("ProducerGroupName"),
	)
	// 启动
	err := p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}

	go func() {
		for i:=0; i<10000; i ++{
			// 发送消息
			body := fmt.Sprintf("this is %d message", i)

			result, err := p.SendSync(context.Background(), &primitive.Message{
				Topic: "testTopic04",
				Body:  []byte(body),
			})
			if err != nil {
				fmt.Printf("send message error: %s\n", err.Error())
			} else {
				fmt.Printf("send message seccess, msg:%s result=%s\n", body, result.String())
			}
			time.Sleep(50*time.Millisecond)
		}
	}()
}

func SubcribeMessage() {
	// 订阅主题、消费
	endPoint := []string{"192.168.10.17:9876"}
	// 创建一个consumer实例
	c, err := rocketmq.NewPushConsumer(consumer.WithNameServer(endPoint),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithGroupName("MyConsumerGroupName"),
	)

	// 订阅topic
	err = c.Subscribe("testTopic04", consumer.MessageSelector{},
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			fmt.Printf("subscribe callback : %v \n", msgs[i])
		}
		return consumer.ConsumeSuccess, nil
	})

	if err != nil {
		fmt.Printf("subscribe message error: %s\n", err.Error())
	}

	// 启动consumer
	err = c.Start()
	if err != nil {
		fmt.Printf("consumer start error: %s\n", err.Error())
		os.Exit(-1)
	}
	//err = c.Shutdown()
	//if err != nil {
	//	fmt.Printf("shutdown Consumer error: %s\n", err.Error())
	//}
}
