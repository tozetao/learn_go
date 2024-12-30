package article

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

type ReadEvent struct {
	Uid       int64 `json:"uid"`
	ArticleID int64 `json:"article_id"`
}

const TopicReadEvent = "article_read"

// Producer 产生各种事件（事件即消息）
type Producer interface {
	// ProduceReadEvent 产生一个用户读取文章的事件
	ProduceReadEvent(event ReadEvent) error
}

type SaramaSyncProducer struct {
	producer sarama.SyncProducer
}

func NewSyncProducer(syncProducer sarama.SyncProducer) Producer {
	return &SaramaSyncProducer{
		producer: syncProducer,
	}
}

func (p *SaramaSyncProducer) ProduceReadEvent(event ReadEvent) error {
	// 把这个事件作为消息投递到kafka中。
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicReadEvent,
		Value: sarama.ByteEncoder(data),
	})
	return err
}
