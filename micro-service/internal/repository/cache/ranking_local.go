package cache

import (
	"context"
	"errors"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"learn_go/webook/internal/domain"
	"time"
)

type LocalCacheRanking struct {
	value *atomicx.Value[[]domain.Article]

	ddl *atomicx.Value[time.Time]

	expiration time.Duration
}

func NewLocalCacheRanking(expiration time.Duration) *LocalCacheRanking {
	return &LocalCacheRanking{
		value:      &atomicx.Value[[]domain.Article]{},
		ddl:        &atomicx.Value[time.Time]{},
		expiration: expiration,
	}
}

func (l *LocalCacheRanking) Set(ctx context.Context, arts []domain.Article) error {
	l.value.Store(arts)
	l.ddl.Store(time.Now().Add(l.expiration))
	return nil
}

func (l *LocalCacheRanking) Get(ctx context.Context) ([]domain.Article, error) {
	arts := l.value.Load()
	ddl := l.ddl.Load()
	if len(arts) == 0 || ddl.Before(time.Now()) {
		return nil, errors.New("本地缓存失效了")
	}
	return arts, nil
}

func (l *LocalCacheRanking) ForceGet(ctx context.Context) ([]domain.Article, error) {
	arts := l.value.Load()
	if len(arts) == 0 {
		return nil, errors.New("本地缓存失效了")
	}
	return arts, nil
}
