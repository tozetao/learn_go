package repository

import (
	"context"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/repository/cache"
)

type RankingRepository interface {
	Get(ctx context.Context) ([]domain.Article, error)

	ReplaceTopN(ctx context.Context, arts []domain.Article) error
}

func (c *cacheRankingRepository) Get(ctx context.Context) ([]domain.Article, error) {
	return c.redisCache.Get(ctx)
}

func (c *cacheRankingRepository) ReplaceTopN(ctx context.Context, articles []domain.Article) error {
	// 考虑只缓存榜单需要的字段
	for i := 0; i < len(articles); i++ {
		articles[i].Content = ""
	}
	return c.redisCache.Set(ctx, articles)
}

func (c *cacheRankingRepository) GetV1(ctx context.Context) ([]domain.Article, error) {
	arts, err := c.localCache.Get(ctx)
	if err == nil {
		return arts, nil
	}
	arts, err = c.redisCache.Get(ctx)
	if err != nil {
		// redis不可用。注：这里没有区分错误，err可能是redis.Nil
		return c.localCache.ForceGet(ctx)
	}
	_ = c.localCache.Set(ctx, arts)
	return arts, nil
}

func (c *cacheRankingRepository) ReplaceTopNV1(ctx context.Context, articles []domain.Article) error {
	// 考虑只缓存榜单需要的字段
	for i := 0; i < len(articles); i++ {
		articles[i].Content = ""
	}
	_ = c.localCache.Set(ctx, articles)
	return c.redisCache.Set(ctx, articles)
}

func NewRankingRepository(redisCache *cache.RedisRanking, localCache *cache.LocalCacheRanking) RankingRepository {
	return &cacheRankingRepository{
		redisCache: redisCache,
		localCache: localCache,
	}
}

type cacheRankingRepository struct {
	//cache cache.RankingCache
	redisCache *cache.RedisRanking
	localCache *cache.LocalCacheRanking
}
