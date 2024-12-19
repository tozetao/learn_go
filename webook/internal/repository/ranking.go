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
	return c.cache.Get(ctx)
}

func (c *cacheRankingRepository) ReplaceTopN(ctx context.Context, articles []domain.Article) error {
	// 考虑只缓存榜单需要的字段
	for i := 0; i < len(articles); i++ {
		articles[i].Content = ""
	}
	return c.cache.Set(ctx, articles)
}

func NewRankingRepository(cache cache.RankingCache) RankingRepository {
	return &cacheRankingRepository{
		cache: cache,
	}
}

type cacheRankingRepository struct {
	cache cache.RankingCache
}
