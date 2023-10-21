package cp_mq

import (
	"github.com/Shopify/sarama"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/nsqio/go-nsq"
)

type RocketMQRelated struct {
	OriMsg		*primitive.MessageExt
	QueueId 	int
	BrokerName 	string
}

type KafkaRelated struct {
	OriMsg		*sarama.ConsumerMessage
	Session 	sarama.ConsumerGroupSession
	Claim		sarama.ConsumerGroupClaim
	Partition 	int32		//分区ID
}

type NsqRelated struct {
	OriMsg		*nsq.Message
	Attempts  	int
}

type Message struct {
	Topic		string		//topic
	Key		string		//key
	Body      	[]byte		//消息内容
	ID        	string 		//源ID
	Timestamp 	int64       	//消息产生时间戮，纳秒
	Offset		int64		//本消息偏移量

	RocketMQ	*RocketMQRelated
	Kafka		*KafkaRelated
	Nsq		*NsqRelated
}
