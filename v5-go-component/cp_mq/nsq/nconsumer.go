package cp_mq_nsq

import (
	"warehouse/v5-go-component/cp_log"
	"context"
	"errors"
	"fmt"
	"runtime"
	"time"

	"warehouse/v5-go-component/cp_mq"
	"github.com/nsqio/go-nsq"
)

type NConsumer struct {
	nc       *nsq.Consumer
	isInit   bool
	conf     *config
	state    cp_mq.ConsumerState
	err      error
	runCount int64 //消费统计
	delayAvg int64 //平均延迟
	delaySum int64 //累计延迟
}

func (this *NConsumer) Init(configStr string, handler cp_mq.ConsumerHandlerFunc) (err error) {
	if this.isInit {
		return
	}

	conf, err := newConfig("consumer", configStr)
	if err != nil {
		err = fmt.Errorf("NSQ：Consumer config Error：%q", err)
		return
	}

	nconf := nsq.NewConfig()
	nc, err := nsq.NewConsumer(conf.Topic, conf.Channel, nconf)
	if err != nil {
		err = fmt.Errorf("NSQ: Consumer Init error %q", err)
		return
	}

	this.nc = nc
	this.conf = conf
	this.isInit = true

	this.nc.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) (err error) {
		startTime := time.Now()
		this.state = cp_mq.CONSUMER_STATE_RUN

		defer recoverConsumerPanic(this)

		// handle the message
		msg := &cp_mq.Message{
			Body:      m.Body,
			ID:        fmt.Sprintf("%s", m.ID),
			Timestamp: m.Timestamp,
			Nsq: 	&cp_mq.NsqRelated {
				OriMsg: m,
				Attempts: int(m.Attempts),
			},
		}

		//todo NSQ是否增加重试机制？
		err, _ = handler(string(msg.Body))
		if err != nil {
			errorInfo := fmt.Sprintf("%s消息处理失败：%s\n消息：%s", this.Name(), err, m.Body)
			cp_log.Error(errorInfo)

			//业务返回错误不处理
			//c.err = err
			//c.Stop()
		}

		this.AddRunCount()
		this.AddDelaySum(time.Now().Sub(startTime).Nanoseconds() / 1000)
		this.state = cp_mq.CONSUMER_STATE_WAIT

		return
	}))

	return
}

func (this *NConsumer) Start() (err error) {
	if !this.isInit {
		err = errors.New(`NSQ：Consumer Start error "must first call the Init func"`)
		return
	}

	if this.state != cp_mq.CONSUMER_STATE_STOP && this.state != cp_mq.CONSUMER_STATE_INIT {
		err = errors.New("NSQ: Consumer is running")
		return
	}

	err = this.nc.ConnectToNSQLookupd(this.conf.Address)
	if err != nil {
		err = fmt.Errorf("NSQ: Consumer ConnectToNSQLookupd error %q", err)
	} else {
		if this.nc.Stats().Connections > 0 {
			this.state = cp_mq.CONSUMER_STATE_WAIT
		} else {
			err = errors.New("NSQ：Consumer connection invalid")
		}
	}

	return
}

func (this *NConsumer) Stop(es ...error) (err error) {
	if !this.isInit {
		err = errors.New(`NSQ：Consumer Stop error "must first call the Init func"`)
		return
	}

	if this.state == cp_mq.CONSUMER_STATE_STOP || this.state == cp_mq.CONSUMER_STATE_STOPING {
		return
	}
	this.state = cp_mq.CONSUMER_STATE_STOPING


	if len(es) == 1 {
		this.err = es[0]
	}
	this.nc.Stop()

	// 阻塞等待消息消费完成
	<-this.nc.StopChan

	this.state = cp_mq.CONSUMER_STATE_STOP
	cp_log.Warning(fmt.Sprintf("NSQ consumer [%s] closed!", this.Name()))

	return
}

func (this *NConsumer) State() cp_mq.ConsumerState {
	return this.state
}

func (this *NConsumer) Error() error {
	return this.err
}

func (this *NConsumer) Name() string {
	return fmt.Sprintf("%s-%s", this.conf.Topic, this.conf.Channel)
}

func (this *NConsumer) Ctx() context.Context {
	return nil
}

func (this *NConsumer) RunCount() int64 {
	return this.runCount
}

func (this *NConsumer) AddRunCount() {
	this.runCount ++
}

func (this *NConsumer) DelayAvg() int64 {
	if this.runCount == 0 {
		return 0
	}
	this.delayAvg = this.delaySum / this.runCount
	return this.delayAvg
}

func (this *NConsumer) DelaySum() int64 {
	return this.delaySum
}

func (this *NConsumer) AddDelaySum(d int64) {
	this.delaySum += d
}

func (this *NConsumer) SetComsumerStatus(s cp_mq.ConsumerState) {
	this.state = s
}

func recoverConsumerPanic(c *NConsumer) {
	if e := recover(); e != nil {
		errorInfo := fmt.Sprintf("%s消息订阅执行故障：%v\n故障堆栈：", c.Name(), e)
		for i := 1; ; i += 1 {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			} else {
				errorInfo += "\n"
			}
			errorInfo += fmt.Sprintf("%v %v", file, line)
		}

		c.Stop(errors.New(errorInfo))
	}
}
