package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"learn_go/webook/internal/domain"
	"log"
	"time"
)

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, user domain.User, expiration time.Duration) error
}

type RedisCache struct {
	client redis.Cmdable
}

func NewUserCache(client redis.Cmdable) UserCache {
	return &RedisCache{
		client: client,
	}
}

func (cache *RedisCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}

func (cache *RedisCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)

	val, err := cache.client.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{}
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		log.Println("json.Unmarshal err:", err)
		return domain.User{}, err
	}
	return user, nil
}

func (cache *RedisCache) Set(ctx context.Context, user domain.User, expiration time.Duration) error {
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return cache.client.Set(ctx, cache.key(user.ID), val, expiration).Err()
}
