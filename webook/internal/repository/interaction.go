package repository

import (
	"context"
	"learn_go/webook/internal/repository/dao"
)

type InteractionRepository interface {
	IncrReadCnt(ctx context.Context, articleID int64) error

	// CreateLikeInfo 创建点赞信息
	CreateLikeInfo(ctx context.Context, articleID int64) error
	// RemoveLikeInfo 移除点赞信息
	RemoveLikeInfo(ctx context.Context, articleID int64) error
}

func (repo interactionRepository) IncrReadCnt(ctx context.Context, articleID int64) error {
	// 1. 数据库自增+1
	err := repo.dao.IncrReadCnt(ctx, articleID)
	if err != nil {
		return err
	}

	// 2. 缓存自增+1
	//    在key存在的情况下，key存在意味着有人访问了该文章，缓存中会载入该文章的交互信息。
}

func (repo interactionRepository) CreateLikeInfo(ctx context.Context, articleID int64) error {
	//TODO implement me
	panic("implement me")
}

func (repo interactionRepository) RemoveLikeInfo(ctx context.Context, articleID int64) error {
	//TODO implement me
	panic("implement me")
}

type interactionRepository struct {
	dao dao.InteractionDao
}

func NewInteractionRepository(dao dao.InteractionDao) InteractionRepository {
	return &interactionRepository{
		dao: dao,
	}
}
