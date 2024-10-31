package ioc

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func NewRedis() redis.Cmdable {
	type RedisConfig struct {
		Addr string
	}

	config := &RedisConfig{}

	err := viper.UnmarshalKey("redis", config)
	if err != nil {
		panic(err)
	}

	fmt.Printf("redis config: %v\n", *config)

	return redis.NewClient(&redis.Options{
		//Addr: "192.168.1.100:6379",
		Addr: config.Addr,
		DB:   0,
	})
}
