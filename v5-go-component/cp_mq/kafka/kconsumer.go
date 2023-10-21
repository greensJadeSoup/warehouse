package cp_mq_kafka

import (
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_mq"
	"warehouse/v5-go-component/cp_util"
	"context"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"runtime"
	"strings"
	"sync"
	"time"
)

type KConsumerItem struct {
	Ready chan bool
	Topics string
	Group string

}

type KConsumer struct {
	cgClient 	sarama.ConsumerGroup
	isInit		bool
	err		error
	conf     	*config
	state    	cp_mq.ConsumerState
	runCount 	int64 //消费统计
	delayAvg 	int64 //平均延迟
	delaySum 	int64 //累计延迟

	ready		chan bool
	ctx		context.Context
	cancel		func()
	wg		sync.WaitGroup			//整个程序关闭的时候，等待所有消费者关闭的同步组

	messageHandler func (body string) (error, cp_constant.MQ_ERR_TYPE)
}

// Setup is run at the beginning of a new session, before ConsumeClaim
// 本函数，在刚启动消费者的时候，会进行本回调。一个消费者只会调用一次。
// 所以处理顺序为：初始化->
//		初始化协程阻塞在ready channel，并调用本函数 ->
//		本函数进行用户自定义处理，并主动关闭 ready channel ->
//		初始化协程阻塞通过ready channel，进入开始消费状态，并阻塞在ctx->Done. 等待主动取消上下文 ->
//		启动多个ConsumeClaim协程，数量根据分区数量而定 ->
//		当遇到主动stop()中调用到ctx->cancel()，主动取消上下文 ->
//		阻塞在stop()中的wg->Wait(),等待所有ConsumeClaim消息处理逻辑处理完毕，退出 ->
//		调用CleanUp()，OfficialRunningCode中wg->Done() ->
//		stop()函数阻塞wg通过
func (this *KConsumer) Setup(session sarama.ConsumerGroupSession) error {
	defer close(this.ready)// Mark the consumer as ready
	this.state = cp_mq.CONSUMER_STATE_RUN
	cp_log.Info(fmt.Sprintf("[Kafka][%s] consumer up and running...", this.conf.Topic))
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
// 本函数Cleanup，需要等到所有ConsumeClaim都退出后（一个消费者有几个分区就有几个ConsumeClaim协程），才会调用。
func (this *KConsumer) Cleanup(session sarama.ConsumerGroupSession) error {
	cp_log.Info(fmt.Sprintf("Kafka consumer [%s] Cleanup", this.conf.Topic))
	return nil
}

// NOTE:
// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Do not move the code below to a goroutine.
// The `ConsumeClaim` itself is called within a goroutine, see:
// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
//
// 本函数ConsumeClaim，当前进程被分配到几个分区，就会启动几个本函数的协程，并且阻塞在下面的range claim.Messages()，等待新消息的到来！！
// 本函数ConsumeClaim，除非主动取消上下文，否则return后还是会重新进来
// session.MarkMessage(): 把offset标记到最新消费的位置，不管前面的offset是否有调用到此方法，都定位到本消息offset
// claim.InitialOffset(): 本次会话从本分区(注意不是本topic)哪个offset开始
// consumerMsg.Offset: 本消息在本分区(注意不是本topic)偏移量
// claim.HighWaterMarkOffset(): 本分区(注意不是本topic)的高度
func (this *KConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	defer recoverConsumerPanic(this, claim)

	cp_log.Info(fmt.Sprintf("[%s]claim.Messages通道进入：%d", claim.Topic(), claim.Partition()))

	for consumerMsg := range claim.Messages() {
		if this.state == cp_mq.CONSUMER_STATE_STOPING || this.state == cp_mq.CONSUMER_STATE_STOP {
			cp_log.Warning("consumer worker协程收到退出信号, 主动退出")
			return nil
		}

		startTime := time.Now()

		msg := &cp_mq.Message {
			ID: "KAFKA" + cp_util.NewGuid(),
			Topic: consumerMsg.Topic,
			Key: string(consumerMsg.Key),
			Body: consumerMsg.Value,
			Offset: consumerMsg.Offset,

			Kafka: &cp_mq.KafkaRelated {
				Session: session,
				Claim: claim,
				OriMsg: consumerMsg,
				Partition: consumerMsg.Partition,
			},
		}

		cp_log.Info(fmt.Sprintf("[%s][%d]new message: %s", this.conf.Topic, claim.Partition(), string(msg.Body)))

		// 正式进入消息处理回调
		err, errType := this.messageHandler(string(msg.Body))
		if err != nil {
			cp_log.Error(fmt.Sprintf("[KafKa][%s]消息处理失败, 类型[%d]: partition=%d, offset=%d, key=%s, error:%s, msgID=%s, message:%s",
				msg.Topic, errType, msg.Kafka.Partition, msg.Offset, msg.Key, err.Error(), msg.ID, string(msg.Body)))
			//todo kafka需要自己实现死信功能！！！
		} else {
			cp_log.Info(fmt.Sprintf("[Kafka][%s]消息消费成功: partition=%d, offset=%d, key=%s, value=%s",
				msg.Topic, msg.Kafka.Partition, msg.Offset, msg.Key, string(msg.Body)))
		}

		session.MarkMessage(msg.Kafka.OriMsg, "")
		this.AddRunCount()
		this.AddDelaySum(time.Now().Sub(startTime).Nanoseconds() / 1000)
	}

	return nil
}

func (this *KConsumer) Init(configStr string, handler cp_mq.ConsumerHandlerFunc) (err error) {
	if this.isInit {
		return
	}

	this.state = cp_mq.CONSUMER_STATE_INIT

	conf, err := newConfig("consumer", configStr)
	if err != nil {
		err = fmt.Errorf("Kafka：Consumer config Error：%q", err)
		return
	}
	cp_log.Info("[Kafka]" + conf.Topic + " initing...")

	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	//创建订阅者群，集群地址发布者代码里已定义
	this.cgClient, err = sarama.NewConsumerGroup(conf.Address, conf.Group, config)
	if err != nil {
		return errors.New(fmt.Sprintf("Error creating consumer group client[%v]: %s", conf.Address, err.Error()))
	}

	this.conf = conf
	this.isInit = true
	this.ready = make(chan bool)
	this.messageHandler = handler

	return
}

func (this *KConsumer) Start() (err error) {
	if !this.isInit {
		err = errors.New(`Kafka: Consumer Start error "must first call the Init func"`)
		return
	}

	if this.state != cp_mq.CONSUMER_STATE_STOP && this.state != cp_mq.CONSUMER_STATE_INIT {
		err = errors.New("Kafka: Consumer is running")
		return
	}

	go this.OfficialRunningCode()

	return nil
}

func (this *KConsumer) Stop(es ...error) (err error) {
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
		cp_log.Warning(fmt.Sprintf("Kafka[%s]Consumer Stop: %s", this.conf.Topic, this.err))
	}

	//this.stop<-true //取消消费者中的重试处理，避免等太久
	this.cancel()	//取消上下文
	this.wg.Wait()

	for i := 0; i < 3; i ++ { //关闭客户端
		if err = this.cgClient.Close(); err == nil {
			break
		}
		cp_log.Error("Error closing client: " + err.Error())
	}

	this.state = cp_mq.CONSUMER_STATE_STOP
	cp_log.Warning(fmt.Sprintf("Kafka consumer [%s] closed success.", this.Name()))

	return
}

func (this *KConsumer) State() cp_mq.ConsumerState {
	return this.state
}

func (this *KConsumer) Error() error {
	return this.err
}

func (this *KConsumer) Ctx() context.Context {
	return this.ctx
}

func (this *KConsumer) Name() string {
	return fmt.Sprintf("%s %s", this.conf.Topic, this.conf.Group)
}

func (this *KConsumer) RunCount() int64 {
	return this.runCount
}

func (this *KConsumer) DelaySum() int64 {
	return this.delaySum
}

func (this *KConsumer) AddRunCount() {
	this.runCount ++
}

func (this *KConsumer) DelayAvg() int64 {
	if this.runCount == 0 {
		return 0
	}
	this.delayAvg = this.delaySum / this.runCount
	return this.delayAvg
}

func (this *KConsumer) AddDelaySum(d int64) {
	this.delaySum += d
}

func (this *KConsumer) SetComsumerStatus(s cp_mq.ConsumerState) {
	this.state = s
}

func recoverConsumerPanic(c *KConsumer, claim sarama.ConsumerGroupClaim) {
	cp_log.Info(fmt.Sprintf("[%s]claim.Messages通道退出：%d", claim.Topic(), claim.Partition()))

	if e := recover(); e != nil {
		errorInfo := fmt.Sprintf("Topic:[%s] Partition:[%d] Kafka消息订阅执行故障：%v\n故障堆栈：", c.Name(), claim.Partition(), e)
		for i := 1; ; i += 1 {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			} else {
				errorInfo += "\n"
			}
			errorInfo += fmt.Sprintf("%v %v", file, line)
		}

		go c.Stop(errors.New(errorInfo))
	}
}

func (this *KConsumer) OfficialRunningCode() {

	this.ctx, this.cancel = context.WithCancel(context.Background())//创建一个上下文对象，实际项目中也一定不要设置超时（当然，按你项目需求，我是没见过有项目需求要多少时间后取消订阅的）
	this.wg.Add(1)//创建同步组

	go func() {
		defer this.wg.Done()
		for {
			/**
			  官方说：`订阅者`应该在无限循环内调用
			  当`发布者`发生变化时
			  需要重新创建`订阅者`会话以获得新的声明

			  所以这里把订阅者放在了循环体内
			*/

			if err := this.cgClient.Consume(this.ctx, strings.Split(this.conf.Topic, ","), this); err != nil {
				cp_log.Error(fmt.Sprintf("[kafka][%s]Error from consumer: %s", this.conf.Topic, err.Error()))
			}
			// 检查上下文是否被取消，收到取消信号应当立刻在本协程中取消循环
			if this.ctx.Err() != nil {
				//被上下文主动取消
				cp_log.Warning(fmt.Sprintf("[kafka][%s]consumer cancel by context.cancel()!", this.conf.Topic))
				return
			}
			//获取订阅者准备就绪信号
			cp_log.Warning(fmt.Sprintf("[kafka][%s]consumer cancel by remote!", this.conf.Topic))
			this.ready = make(chan bool)
		}
	}()

	<-this.ready //等待SetUp()执行完
	//TODO 可以做一些初始化完毕后，自定义的处理, 但是要在一个协程中写个死循环

	//<-this.ctx.Done() //开始消费后，阻塞在这里，直到主动取消消费(取消上下文)
	//cp_log.Info("ctx.Done()!")
	//this.Stop(errors.New(fmt.Sprintf("[Kafka][%s] context cancelled.", this.conf.Topics)))

	return
}

