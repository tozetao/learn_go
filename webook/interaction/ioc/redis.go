package ioc

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"learn_go/webook/pkg/logger"
)

func NewRedis(log logger.LoggerV2) redis.Cmdable {
	type RedisConfig struct {
		Addr string
	}

	config := &RedisConfig{}

	err := viper.UnmarshalKey("redis", config)
	if err != nil {
		panic(err)
	}

	log.Info("", logger.Field{Key: "redis_config", Value: config})

	return redis.NewClient(&redis.Options{
		//Addr: "192.168.1.100:6379",
		Addr: config.Addr,
		DB:   0,
	})
}
