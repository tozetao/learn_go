package repository

import (
	"context"
	"learn_go/webook/internal/repository/cache"
)

var (
	ErrTooManySend   = cache.ErrTooManySend
	ErrTooManyVerify = cache.ErrTooManyVerify
)

type CodeRepository struct {
	codeCache cache.CodeCache
}

func NewCodeRepository(codeCache cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		codeCache: codeCache,
	}
}

func (repo *CodeRepository) Store(ctx context.Context, biz string, phone string, code string) error {
	return repo.codeCache.Set(ctx, biz, phone, code)
}

func (repo *CodeRepository) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return repo.codeCache.Verify(ctx, biz, phone, inputCode)
}
