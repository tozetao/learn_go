package ioc

import (
	"github.com/redis/go-redis/v9"
	"learn_go/webook/config"
)

func NewRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		//Addr: "192.168.1.100:6379",
		Addr: config.Config.Redis.Addr,
		DB:   0,
	})
}
