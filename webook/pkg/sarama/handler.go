package sarama

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"learn_go/webook/pkg/logger"
)

// 消费者。启动消费者，消费消息。
// 我需要一个kafka客户端，然后开始消费消息?

// Handler sarama.ConsumerGroupHandler的实现
type Handler[T any] struct {
	l  logger.LoggerV2
	fn func(event T, topic string, partition int32, offset int64) error
}

func NewHandler[T any](l logger.LoggerV2, fn func(event T, topic string, partition int32, offset int64) error) sarama.ConsumerGroupHandler {
	return &Handler[T]{
		l:  l,
		fn: fn,
	}
}

func (c *Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// c.l.Info(fmt.Sprintf("topic: %s, partition: %d", claim.Topic(), claim.Partition()))
	messages := claim.Messages()
	for msg := range messages {
		// c.l.Info(fmt.Sprintf("topic: %s, partition: %d, Offset: %d", msg.Topic, msg.Partition, msg.Offset))
		var data T
		err := json.Unmarshal(msg.Value, &data)
		// 格式解析错误，这意味着生产者投递的消息格式有问题，不需要重试。
		if err != nil {
			continue
		}

		//err = c.fn(event, msg.Topic, msg.Partition, msg.Offset)
		//if err != nil {
		//	// 记录日志
		//	// 重试：使用装饰器来实现重试的逻辑
		//}
		for i := 0; i < 3; i++ {
			err = c.fn(data, msg.Topic, msg.Partition, msg.Offset)
			if err == nil {
				break
			}
			c.l.Error("消费消息失败", logger.Error(err),
				logger.String("topic", msg.Topic),
				logger.Int32("partition", msg.Partition),
				logger.Int64("offset", msg.Offset))
		}

		if err != nil {
			c.l.Error("重试失败", logger.Error(err),
				logger.String("topic", msg.Topic),
				logger.Int32("partition", msg.Partition),
				logger.Int64("offset", msg.Offset))
		}
		// 无论成功或失败都提交消息
		session.MarkMessage(msg, "")
	}
	return nil
}
