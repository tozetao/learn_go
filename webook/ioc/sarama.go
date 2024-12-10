package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"learn_go/webook/pkg/logger"
)

// NewSaramaClient sarama client
func NewSaramaClient(log logger.LoggerV2) sarama.Client {
	// 可以通过读取配置来进行初始化
	type Config struct {
		Addr []string
	}
	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	saramaCfg := sarama.NewConfig()
	client, err := sarama.NewClient(config.Addr, &saramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}
