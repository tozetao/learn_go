package repository

import (
	"context"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/repository/cache"
	"learn_go/webook/internal/repository/dao"
)

type InteractionRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizID int64) error

	// IncrLike 创建点赞信息
	IncrLike(ctx context.Context, uid int64, biz string, bizID int64) error
	// DecrLike 移除点赞信息
	DecrLike(ctx context.Context, uid int64, biz string, bizID int64) error
	AddFavorite(ctx context.Context, uid int64, favoriteID int64, biz string, bizID int64) error
}

func (repo *interactionRepository) AddFavorite(ctx context.Context, uid int64, favoriteID int64, biz string, bizID int64) error {
	err := repo.dao.InsertFavorite(ctx, dao.UserFavorite{
		Uid:        uid,
		FavoriteID: favoriteID,
		Biz:        biz,
		BizID:      bizID,
	})
	if err != nil {
		return err
	}
	// 增加缓存计数
	err = repo.cache.IncrFavoriteCnt(ctx, biz, bizID)
	if err != nil {
		// 记录日志
	}
	return nil
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

func (repo *interactionRepository) IncrReadCnt(ctx context.Context, biz string, bizID int64) error {
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

func (repo *interactionRepository) IncrLike(ctx context.Context, uid int64, biz string, bizID int64) error {
	// 1. 插入我点赞的文章
	// 2. 文章点赞+1

	err := repo.dao.InsertLikeInfo(ctx, uid, biz, bizID)
	if err != nil {
		return err
	}

	// 缓存计数+1
	return repo.cache.IncrLikeCnt(ctx, biz, bizID)
}

func (repo *interactionRepository) DecrLike(ctx context.Context, uid int64, biz string, bizID int64) error {

	err := repo.dao.DeleteLikeInfo(ctx, uid, biz, bizID)
	if err != nil {
		return err
	}

	// 缓存计数-1
	return repo.cache.DecrLikeCnt(ctx, biz, bizID)
}

func (repo *interactionRepository) toEntity(inter domain.Interaction) dao.Interaction {
	return dao.Interaction{
		ID:    inter.ID,
		Biz:   inter.Biz,
		BizID: inter.ID,

		Favorites: inter.Favorites,
		ReadCnt:   inter.Views,
		Likes:     inter.Likes,

		CTime: inter.CTime.UnixMilli(),
		UTime: inter.UTime.UnixMilli(),
	}
}
