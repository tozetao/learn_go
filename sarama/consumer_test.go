package sarama_test

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	//config := sarama.NewConfig()
	//consumer, err := sarama.NewConsumerGroup(addr, "demo", config)
	//require.NoError(t, err)
	//
	//// 阻塞在这里?
	//err = consumer.Consume(context.Background(), []string{"my_topic"}, ConsumerHandler{})
	//t.Logf("close consumer, %v\n", err)
	//
	//// 生产者关闭服务
	//require.NoError(t, err)

	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumerGroup([]string{"localhost:9094"}, "demo", cfg)
	assert.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()
	start := time.Now()
	err = consumer.Consume(ctx,
		[]string{"my_topic"}, &ConsumerHandler{t})
	assert.NoError(t, err)
	t.Log(time.Since(start).String())
}

type ConsumerHandler struct {
	t *testing.T
}

func (c *ConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	c.t.Logf("Setup, claims: %v\n", session.Claims())
	return nil
}

func (c *ConsumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	c.t.Log("Cleanup\n")
	return nil
}

func (c *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	c.t.Logf("topic: %s, partition: %d", claim.Topic(), claim.Partition())
	msgs := claim.Messages()

	const batchSize = 1000

	for {
		var eg errgroup.Group
		done := false
		msgContainer := make([]*sarama.ConsumerMessage, 0, batchSize)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)

		for i := 0; i < batchSize; i++ {
			if done {
				break
			}
			select {
			case <-ctx.Done():
				done = true
			case msg, ok := <-msgs:
				if !ok {
					cancel()
					return nil
				}
				msgContainer = append(msgContainer, msg)
				// 并发处理消息
				eg.Go(func() error {
					c.t.Log(string(msg.Value))
					return nil
				})
			}
		}

		cancel()
		err := eg.Wait()
		if err != nil {
			// 记录错误
		}

		// 批量提交消息
		for _, msg := range msgContainer {
			session.MarkMessage(msg, "")
		}
	}
}

func (c *ConsumerHandler) ConsumeClaimV1(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	c.t.Logf("topic: %s, partition: %d", claim.Topic(), claim.Partition())
	msgs := claim.Messages()
	for msg := range msgs {
		c.t.Logf("message: %s, topic: %s, partition: %d, offset: %d", string(msg.Value), msg.Topic, msg.Partition, msg.Offset)
		c.t.Log(msg.Offset % 2)
		if msg.Offset%2 == 0 {
			session.MarkMessage(msg, "")
		}
	}
	return nil
}
