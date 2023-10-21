package cp_mq_kafka

import (
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"time"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_mq"
)

type KProducer struct {
	pClient		sarama.SyncProducer
	isInit		bool
	conf     	*config
	state    	cp_mq.ProducerState
	runCount 	int64 //发布统计
	delayAvg 	int64 //平均延迟
	delaySum 	int64 //累计延迟

	messageHandler func ([]byte) error
}

func (this *KProducer) Init(configStr string) (err error) {
	if this.isInit {
		return
	}

	conf, err := newConfig("producer", configStr)
	if err != nil {
		err = fmt.Errorf("Kafka：Consumer config Error：%q", err)
		return
	}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second
	config.Producer.RequiredAcks = -1 // 设置可靠性，类似于MQTT QOS

	this.pClient, err = sarama.NewSyncProducer(conf.Address, config)
	if err != nil {
		log.Printf("sarama.NewSyncProducer err, message=%s \n", err)
		return
	}

	this.conf = conf
	this.isInit = true
	this.state = cp_mq.PRODUCER_STATE_WAIT

	return
}

func (this *KProducer) Publish(value []byte, key string) (err error) {
	if !this.isInit {
		return errors.New(`Kafka：Publish error "must first call the Init func"`)
	}

	if this.state == cp_mq.PRODUCER_STATE_STOP || this.state == cp_mq.PRODUCER_STATE_STOPING {
		err = errors.New(`Kafka：Publish error "producer has stoped."`)
		return
	}

	startTime := time.Now()
	this.state = cp_mq.PRODUCER_STATE_RUN

	msg := &sarama.ProducerMessage {
		Topic: this.conf.Topic,
		Value: sarama.ByteEncoder(value),
	}

	if key != "" {
		msg.Key = sarama.StringEncoder(key)
	}

	partition, offset, err := this.pClient.SendMessage(msg)
	if err != nil {
		err = errors.New("Kafka: send message err = " + err.Error())
		return
	}

	this.AddRunCount()
	this.AddDelaySum(time.Now().Sub(startTime).Nanoseconds() / 1000)

	this.state = cp_mq.PRODUCER_STATE_WAIT
	cp_log.Debug(fmt.Sprintf("[Kafka]消息发送成功，partition=%d, offset=%d \n", partition, offset))

	return nil
}

func (this *KProducer) Stop() (err error) {
	if !this.isInit {
		err = errors.New(`Kafka：Producer Stop error "must first call the Init func"`)
		return
	}

	this.state = cp_mq.PRODUCER_STATE_STOPING
	//todo 这里可能要等所有生产者退出后，才能关闭客户端，以后完善

	for i := 0; i < 3; i ++ { //关闭客户端
		if err = this.pClient.Close(); err == nil {
			break
		}
		cp_log.Error("Error closing client: " + err.Error())
	}

	this.state = cp_mq.PRODUCER_STATE_STOP
	cp_log.Warning(fmt.Sprintf("Kafka publicer [%s] closed success.", this.Name()))

	return
}

func (this *KProducer) State() cp_mq.ProducerState {
	return this.state
}

func (this *KProducer) Name() string {
	return fmt.Sprintf("%s", this.conf.Topic)
}

func (this *KProducer) RunCount() int64 {
	return this.runCount
}

func (this *KProducer) DelaySum() int64 {
	return this.delaySum
}

func (this *KProducer) AddRunCount() {
	this.runCount ++
}

func (this *KProducer) DelayAvg() int64 {
	if this.runCount == 0 {
		return 0
	}
	this.delayAvg = this.delaySum / this.runCount
	return this.delayAvg
}

func (this *KProducer) AddDelaySum(d int64) {
	this.delaySum += d
}