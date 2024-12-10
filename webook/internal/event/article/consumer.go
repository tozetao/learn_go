package article

import (
	"github.com/IBM/sarama"
	"learn_go/webook/internal/event"
	"learn_go/webook/pkg/logger"
)

// 阅读文章事件的消费者

// GroupRead Group ID
const GroupRead = "group:article_read"

type Consumer struct {
	cg sarama.ConsumerGroup

	l logger.LoggerV2
}

func NewConsumer(cg sarama.ConsumerGroup, l logger.LoggerV2) event.Consumer {
	return &Consumer{
		cg: cg,
		l:  l,
	}
}

func (c *Consumer) Start() error {

	// 创建多个消费者组来处理文章阅读事件

	// sarama.NewConsumerGroupFromClient()

	//c.cg.Consume(context.Background(),
	//	[]string{TopicReadEvent},
	//	saramax.NewHandler[ReadEvent](c.l, c.Consume))
	return nil
}

func (c *Consumer) Consume(event ReadEvent, topic string, partition int32, offset int64) error {
	// 文章读取消费

	return nil
}
