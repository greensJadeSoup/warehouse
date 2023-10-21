package cp_mq_rmq

import (
	"context"
	"errors"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"time"
	"warehouse/v5-go-component/cp_log"

	"warehouse/v5-go-component/cp_mq"
)

type RProducer struct {
	rp	 rocketmq.Producer
	conf     *config
	state    cp_mq.ProducerState
	isInit   bool
	runCount int64 //消费统计
	delayAvg int64 //平均延迟
	delaySum int64 //累计延迟
}

func (this *RProducer) Init(configStr string) (err error) {
	this.conf, err = newConfig("producer", configStr)
	if err != nil {
		err = fmt.Errorf("Producer config Error：%q", err)
		return
	}

	// 创建一个producer实例
	rp, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver(this.conf.Address)),
		producer.WithRetry(2),
	)
	if err != nil {
		err = fmt.Errorf("RocketMQ: Producer Create Error：%q", err)
		return
	}

	// 启动
	err = rp.Start()
	if err != nil {
		err = fmt.Errorf("RocketMQ: Start Producer Error: %s", err.Error())
		return
	}

	this.rp = rp
	this.isInit = true
	this.state = cp_mq.PRODUCER_STATE_WAIT

	return
}

func (this *RProducer) Publish(msg []byte, key string) (err error) {
	if !this.isInit {
		err = errors.New(`RocketMQ：Publish error "must first call the Init func"`)
		return
	}

	if this.state == cp_mq.PRODUCER_STATE_STOP {
		err = errors.New(`RocketMQ：Publish error "producer has stoped."`)
		return
	}

	startTime := time.Now()
	this.state = cp_mq.PRODUCER_STATE_RUN

	result, err := this.rp.SendSync(
		context.Background(),
		&primitive.Message {
			Topic: this.conf.Topic,
			Body:  msg,
	})
	if err != nil {
		err = fmt.Errorf("RocketMQ: Publish error %q", err)
		cp_log.Error(err.Error())
	} else {
		cp_log.Info(fmt.Sprintf("send message seccess, msg:%s result=%s\n", msg, result.String()))
	}

	this.AddRunCount()
	this.AddDelaySum(time.Now().Sub(startTime).Nanoseconds() / 1000)
	this.state = cp_mq.PRODUCER_STATE_WAIT

	return
}

func (this *RProducer) Stop() (err error) {
	if !this.isInit {
		err = errors.New(`RocketMQ：Producer Stop error "must first call the Init func"`)
		return
	}

	this.state = cp_mq.PRODUCER_STATE_STOP
	this.rp.Shutdown()
	return
}

func (this *RProducer) State() cp_mq.ProducerState {
	return this.state
}

func (this *RProducer) Name() string {
	return this.conf.Topic
}

func (this *RProducer) RunCount() int64 {
	return this.runCount
}

func (this *RProducer) DelaySum() int64 {
	return this.delaySum
}

func (this *RProducer) AddRunCount() {
	this.runCount ++
}

func (this *RProducer) DelayAvg() int64 {
	if this.runCount == 0 {
		return 0
	}
	this.delayAvg = this.delaySum / this.runCount
	return this.delayAvg
}

func (this *RProducer) AddDelaySum(d int64) {
	this.delaySum += d
}