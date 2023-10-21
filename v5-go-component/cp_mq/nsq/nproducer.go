package cp_mq_nsq

import (
	"errors"
	"fmt"
	"time"

	"warehouse/v5-go-component/cp_mq"
	"github.com/nsqio/go-nsq"
)

type NProducer struct {
	np       *nsq.Producer
	isInit   bool
	conf     *config
	state    cp_mq.ProducerState
	runCount int64 //消费统计
	delayAvg int64 //平均延迟
	delaySum int64 //累计延迟
}

func (this *NProducer) Init(configStr string) (err error) {
	conf, err := newConfig("producer", configStr)
	if err != nil {
		err = fmt.Errorf("Producer config Error：%q", err)
		return
	}

	conf.Address, err = getNsqdService(conf)
	if err != nil {
		return
	}

	nconf := nsq.NewConfig()
	np, err := nsq.NewProducer(conf.Address, nconf)
	if err != nil {
		err = fmt.Errorf("NSQ: Producer Init error %q", err)
		return
	}

	producerLogger := &logger{}

	np.SetLogger(producerLogger, nsq.LogLevelWarning)
	this.np = np
	this.conf = conf
	this.isInit = true

	this.state = cp_mq.PRODUCER_STATE_WAIT
	return
}

func (this *NProducer) Publish(msg []byte, key string) (err error) {
	if !this.isInit {
		err = errors.New(`NSQ：Publish error "must first call the Init func"`)
		return
	}

	if this.state == cp_mq.PRODUCER_STATE_STOP {
		err = errors.New(`NSQ：Publish error "producer has stoped."`)
		return
	}

	startTime := time.Now()
	this.state = cp_mq.PRODUCER_STATE_RUN

	err = this.np.Publish(this.conf.Topic, msg)
	if err != nil {
		err = fmt.Errorf("NSQ: Publish error %q", err)
	}
	delay := time.Now().Sub(startTime).Nanoseconds() / 1000
	this.runCount += 1
	this.delaySum += delay
	this.delayAvg = this.delaySum / this.runCount

	this.state = cp_mq.PRODUCER_STATE_WAIT

	return
}

func (this *NProducer) Stop() (err error) {
	if !this.isInit {
		err = errors.New(`NSQ：Producer Stop error "must first call the Init func"`)
		return
	}

	this.state = cp_mq.PRODUCER_STATE_STOP
	this.np.Stop()
	return
}

func (this *NProducer) State() cp_mq.ProducerState {
	return this.state
}

func (this *NProducer) Name() string {
	return this.conf.Topic
}

func (this *NProducer) RunCount() int64 {
	return this.runCount
}

func (this *NProducer) DelaySum() int64 {
	return this.delaySum
}

func (this *NProducer) AddRunCount() {
	this.runCount ++
}

func (this *NProducer) DelayAvg() int64 {
	if this.runCount == 0 {
		return 0
	}
	this.delayAvg = this.delaySum / this.runCount
	return this.delayAvg
}

func (this *NProducer) AddDelaySum(d int64) {
	this.delaySum += d
}