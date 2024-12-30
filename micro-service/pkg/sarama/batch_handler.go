package sarama

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"learn_go/webook/pkg/logger"
	"time"
)

// BatchHandler 批量处理sarama.ConsumerGroupHandler的消息
type BatchHandler[T any] struct {
	l             logger.LoggerV2
	fn            func(message []*sarama.ConsumerMessage, events []T) error
	batchDuration time.Duration
	batchSize     int
}

func NewBatchHandler[T any](l logger.LoggerV2, fn func(message []*sarama.ConsumerMessage, events []T) error) sarama.ConsumerGroupHandler {
	return &BatchHandler[T]{
		l:             l,
		fn:            fn,
		batchSize:     256,
		batchDuration: time.Second,
	}
}

func (c *BatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *BatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgCh := claim.Messages()

	for {
		// 1. 在一段时间内获取一批消息。
		// 2. 处理超时但是未获取足够的消息的情况
		//var lastMsg *sarama.ConsumerMessage

		done := false
		ctx, cancel := context.WithTimeout(context.Background(), c.batchDuration)

		// 存储解码的消息
		ts := make([]T, 0, c.batchSize)
		// 存储消息本身
		messages := make([]*sarama.ConsumerMessage, 0, c.batchSize)

		//c.l.Info("进入循环")
		for i := 0; i < c.batchSize && !done; i++ {
			select {
			case msg, ok := <-msgCh:
				if !ok {
					cancel()
					return nil
				}
				//lastMsg = msg

				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					// 解码错误只能记录日志，处理下一条消息
					c.l.Error("消息解码失败", logger.Error(err),
						logger.String("topic", msg.Topic),
						logger.Int32("partition", msg.Partition),
						logger.Int64("offset", msg.Offset))
					continue
				}
				ts = append(ts, t)
				messages = append(messages, msg)
			case <-ctx.Done():
				done = true
			}
		}
		//c.l.Info("推出循环", logger.String("消息数量: ", strconv.Itoa(len(messages))))

		cancel()
		if len(messages) == 0 {
			continue
		}

		// 把消息批量传递给调用方
		err := c.fn(messages, ts)
		if err != nil {
			// 记录错误
			c.l.Error("业务方调用失败", logger.Error(err))
		}
		c.l.Info(fmt.Sprintf("处理的消息数量: %v", len(messages)))

		// 只提交最后一条消息，其他分片的消息未提交怎么办? 还是要批量提交
		//session.MarkMessage(lastMsg, "")
		for _, m := range messages {
			session.MarkMessage(m, "")
		}
	}
}
