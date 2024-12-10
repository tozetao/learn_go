package article

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

type ReadEvent struct {
	Uid       int64
	ArticleID int64
}

const TopicReadEvent = "article_read"

type Producer interface {
	// 产生一个用户读取文章的事件
	produceReadEvent(event ReadEvent) error
}

type SaramaSyncProducer struct {
	producer sarama.SyncProducer
}

func NewSyncProducer(syncProducer sarama.SyncProducer) Producer {
	return &SaramaSyncProducer{
		producer: syncProducer,
	}
}

func (p *SaramaSyncProducer) produceReadEvent(event ReadEvent) error {
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
