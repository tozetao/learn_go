package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func NewSaramaConfig() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.Retry.Max = 3
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	return cfg
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
