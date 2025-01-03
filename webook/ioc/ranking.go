package ioc

import (
	"github.com/redis/go-redis/v9"
	"learn_go/webook/internal/repository/cache"
	"time"
)

func NewRedisRanking(cmd redis.Cmdable) *cache.RedisRanking {
	return cache.NewRedisRanking(cmd, time.Second*10)
}

func NewLocalCacheRanking() *cache.LocalCacheRanking {
	return cache.NewLocalCacheRanking(time.Second * 10)
}
