package cp_mq_rmq

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"sync"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_log"
	"context"
	"errors"
	"fmt"
	"runtime"
	"time"

	"warehouse/v5-go-component/cp_mq"
)

type RConsumer struct {
	rc       	rocketmq.PushConsumer
	isInit   	bool
	conf     	*config
	state    	cp_mq.ConsumerState
	err      	error
	wg		sync.WaitGroup	//整个程序关闭的时候，等待所有消费者关闭的同步组

	runCount 	int64 //消费统计
	delayAvg 	int64 //平均延迟
	delaySum 	int64 //累计延迟
}

//消费者说明：
//1.当队列queue数目为1的时候, 且设置了顺序消费consumer.WithConsumerOrder,
//   1.1 多个进程同个消费组, 只有一个进程能占用队列进行消费
//   1.2 多个进程不同消费组, 则多个进程自己消费自己的组的消息
//
//
//
func (this *RConsumer) Init(configStr string, handler cp_mq.ConsumerHandlerFunc) (err error) {
	if this.isInit {
		return
	}

	this.conf, err = newConfig("consumer", configStr)
	if err != nil {
		err = fmt.Errorf("RocketMQ：Consumer config Error：%q", err)
		return
	}

	fmt.Println(this.conf.Address)

	// 创建一个consumer实例
	rc, err := rocketmq.NewPushConsumer(
		consumer.WithNsResolver(primitive.NewPassthroughResolver(this.conf.Address)),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithGroupName(this.conf.Group),
		consumer.WithMaxReconsumeTimes(5),		//消息最大重试消费次数
		consumer.WithConsumeMessageBatchMaxSize(int(1)), 	//批量消费最大推送数量
		consumer.WithConsumerOrder(true),			//顺序消费 这个选项必须加，否则是无序的且会开很多协程
	)

	// 订阅topic
	err = rc.Subscribe(this.conf.Topic, consumer.MessageSelector{},
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			if this.state == cp_mq.CONSUMER_STATE_STOPING || this.state == cp_mq.CONSUMER_STATE_STOP {
				cp_log.Warning("RocketMQ: consumer worker协程收到退出信号, 主动退出")
				return consumer.Rollback, err
			}

			this.wg.Add(1)
			defer this.wg.Done()

			for i := range msgs {
				startTime := time.Now()
				this.state = cp_mq.CONSUMER_STATE_RUN

				defer recoverConsumerPanic(this)

				msg := &cp_mq.Message{
					Topic:		msgs[i].Topic,
					Key: 		msgs[i].GetKeys(),
					Body:      	msgs[i].Body,
					ID:        	fmt.Sprintf("%s", msgs[i].MsgId),
					Timestamp: 	msgs[i].BornTimestamp,
					Offset: 	msgs[i].QueueOffset,

					RocketMQ: 	&cp_mq.RocketMQRelated {
						OriMsg: 	msgs[i],
						QueueId: 	msgs[i].Queue.QueueId,
						BrokerName: 	msgs[i].Queue.BrokerName,
					},
				}

				err, consumeResult := handler(string(msg.Body))
				if err != nil {
					cp_log.Error(fmt.Sprintf("[RocketMQ][%s]消息处理失败, 类型[%d]: queueId=%d, offset=%d, key=%s, error:%s, msgID=%s, message:%s",
						msg.Topic, consumeResult, msg.RocketMQ.QueueId, msg.Offset, msg.Key, err.Error(), msg.ID, string(msg.Body)))

					if consumeResult == cp_constant.MQ_ERR_TYPE_RECOVERABLE { //可恢复的错误
						return consumer.ConsumeRetryLater, err
					} else if consumeResult == cp_constant.MQ_ERR_TYPE_UNRECOVERABLE { //不可恢复的错误
						return consumer.SuspendCurrentQueueAMoment, err
					}
				} else {
					cp_log.Info(fmt.Sprintf("[RocketMQ][%s]消息消费成功: queueId=%d, offset=%d, key=%s, value=%s",
						msg.Topic, msg.RocketMQ.QueueId, msg.Offset, msg.Key, string(msg.Body)))
				}

				this.AddRunCount()
				this.AddDelaySum(time.Now().Sub(startTime).Nanoseconds() / 1000)
			}

			if this.state != cp_mq.CONSUMER_STATE_STOPING && this.state != cp_mq.CONSUMER_STATE_STOP {
				this.state = cp_mq.CONSUMER_STATE_WAIT
			}

			return consumer.ConsumeSuccess, nil
		})

	if err != nil {
		err = fmt.Errorf("RocketMQ: subscribe message error: %s\n", err.Error())
		return
	}

	this.rc = rc
	this.isInit = true

	return
}

func (this *RConsumer) Start() (err error) {
	if !this.isInit {
		err = errors.New(`RocketMQ：Consumer Start error "must first call the Init func"`)
		return
	}

	if this.state != cp_mq.CONSUMER_STATE_STOP && this.state != cp_mq.CONSUMER_STATE_INIT {
		err = errors.New("RocketMQ: Consumer is running")
		return
	}

	// 启动consumer
	err = this.rc.Start()
	if err != nil {
		err = fmt.Errorf("RocketMQ: consumer start error: %s\n", err.Error())
		return
	}

	this.state = cp_mq.CONSUMER_STATE_WAIT

	return
}

func (this *RConsumer) Stop(es ...error) (err error) {
	if !this.isInit {
		err = errors.New(`RocketMQ：Consumer Stop error "must first call the Init func"`)
		return
	}

	if this.state == cp_mq.CONSUMER_STATE_STOP || this.state == cp_mq.CONSUMER_STATE_STOPING {
		return
	}

	this.state = cp_mq.CONSUMER_STATE_STOPING

	if len(es) == 1 {
		this.err = es[0]
	}

	cp_log.Info(fmt.Sprintf("[RocketMQ]Topic:%s ready to shutdown", this.Name()))
	this.wg.Wait()	// 阻塞等待消息消费完成
	this.rc.Shutdown()
	cp_log.Info(fmt.Sprintf("[RocketMQ]Topic:%s shutdown success", this.Name()))

	this.state = cp_mq.CONSUMER_STATE_STOP

	return
}

func (this *RConsumer) State() cp_mq.ConsumerState {
	return this.state
}

func (this *RConsumer) Error() error {
	return this.err
}

func (this *RConsumer) Name() string {
	return fmt.Sprintf("%s-%s", this.conf.Topic, this.conf.Group)
}

func (this *RConsumer) Ctx() context.Context {
	return nil
}

func (this *RConsumer) RunCount() int64 {
	return this.runCount
}

func (this *RConsumer) AddRunCount() {
	this.runCount ++
}

func (this *RConsumer) DelayAvg() int64 {
	if this.runCount == 0 {
		return 0
	}
	this.delayAvg = this.delaySum / this.runCount
	return this.delayAvg
}

func (this *RConsumer) DelaySum() int64 {
	return this.delaySum
}

func (this *RConsumer) AddDelaySum(d int64) {
	this.delaySum += d
}

func (this *RConsumer) SetComsumerStatus(s cp_mq.ConsumerState) {
	this.state = s
}

func recoverConsumerPanic(c *RConsumer) {
	if e := recover(); e != nil {
		errorInfo := fmt.Sprintf("%s消息订阅执行故障:%v\n 故障堆栈:", c.Name(), e)
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
