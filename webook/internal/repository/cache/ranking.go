package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"learn_go/webook/internal/domain"
	"time"
)

type RankingCache interface {
	Set(ctx context.Context, arts []domain.Article) error

	Get(ctx context.Context) ([]domain.Article, error)
}

type RedisRanking struct {
	redis      redis.Cmdable
	expiration time.Duration
}

func NewRedisRanking(cache redis.Cmdable, expiration time.Duration) *RedisRanking {
	return &RedisRanking{
		redis:      cache,
		expiration: expiration,
	}
}

func (c *RedisRanking) key() string {
	return "ranking:top_n"
}

func (c *RedisRanking) Set(ctx context.Context, articles []domain.Article) error {
	// 这里可以按 id => article 预加载单篇文章数据。分布式环境下，可以考虑给其他进程发送通知，让其他实例也缓存数据。
	data, err := json.Marshal(articles)
	if err != nil {
		return err
	}
	return c.redis.Set(ctx, c.key(), data, c.expiration).Err()
}

func (c *RedisRanking) Get(ctx context.Context) ([]domain.Article, error) {
	data, err := c.redis.Get(ctx, c.key()).Bytes()
	if err != nil {
		return nil, err
	}
	var arts []domain.Article
	err = json.Unmarshal(data, &arts)
	if err != nil {
		return nil, err
	}
	return arts, nil
}
