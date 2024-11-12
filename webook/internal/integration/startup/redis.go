package startup

import "github.com/redis/go-redis/v9"

func NewRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr: "192.168.1.100:6379",
		DB:   0,
	})
}
