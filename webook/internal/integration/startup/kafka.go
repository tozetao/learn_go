package startup

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	event2 "learn_go/webook/interaction/event/article"
	"learn_go/webook/internal/event"
)

func NewSaramaConfig() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.Retry.Max = 3
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	return cfg
}

// NewConsumerClient 构建消息队列的消费者客户端
func NewConsumerClient(saramaCfg *sarama.Config) sarama.Client {
	// 可以通过读取配置来进行初始化
	type Config struct {
		Addrs []string
	}
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}

	client, err := sarama.NewClient(cfg.Addrs, saramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}

func NewConsumers(articleEventConsumer *event2.BatchReadEventConsumer) []event.Consumer {
	return []event.Consumer{articleEventConsumer}
}

func NewSyncProducer(saramaCfg *sarama.Config) sarama.SyncProducer {
	// 可以通过读取配置来进行初始化
	type Config struct {
		Addrs []string
	}
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}

	producer, err := sarama.NewSyncProducer(cfg.Addrs, saramaCfg)
	if err != nil {
		panic(err)
	}
	return producer
}
