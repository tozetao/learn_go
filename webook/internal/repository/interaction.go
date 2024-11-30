package repository

import (
	"context"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/repository/cache"
	"learn_go/webook/internal/repository/dao"
)

type InteractionRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizID int64) error

	// CreateLikeInfo 创建点赞信息
	CreateLikeInfo(ctx context.Context, uid int64, interaction domain.Interaction) error
	// RemoveLikeInfo 移除点赞信息
	RemoveLikeInfo(ctx context.Context, articleID int64) error
}

func (repo interactionRepository) IncrReadCnt(ctx context.Context, biz string, bizID int64) error {
	// 1. 数据库自增+1
	err := repo.dao.IncrReadCnt(ctx, biz, bizID)
	if err != nil {
		return err
	}

	// 2. 缓存自增+1。
	err = repo.cache.IncrReadCnt(ctx, biz, bizID)
	if err != nil {
		// 记录日志错误
	}
	return nil
}

func (repo interactionRepository) CreateLikeInfo(ctx context.Context, uid int64, interaction domain.Interaction) error {
	// 1. 插入我点赞的文章
	// 2. 文章点赞+1
	return nil
}

func (repo interactionRepository) RemoveLikeInfo(ctx context.Context, articleID int64) error {
	//TODO implement me
	panic("implement me")
}

type interactionRepository struct {
	dao   dao.InteractionDao
	cache cache.InteractionCache
}

func NewInteractionRepository(dao dao.InteractionDao, cache cache.InteractionCache) InteractionRepository {
	return &interactionRepository{
		dao:   dao,
		cache: cache,
	}
}
