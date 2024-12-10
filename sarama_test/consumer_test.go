package sarama_test

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
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
	for msg := range msgs {
		c.t.Logf("message: %s, topic: %s, partition: %d, offset: %d", string(msg.Value), msg.Topic, msg.Partition, msg.Offset)

		if msg.Offset%2 == 0 {
			// 忽略offset为偶数的消息
			c.t.Log("ignore the message.")
			continue
		} else {
			session.MarkMessage(msg, "")
		}
	}
	return nil
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
