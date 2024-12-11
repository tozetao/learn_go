package repository

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/repository/cache"
	"learn_go/webook/internal/repository/dao"
	"time"
)

type InteractionRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizID int64) error

	BatchIncrReadCnt(ctx context.Context, bizs []string, bizIDs []int64) error

	// IncrLike 创建点赞信息
	IncrLike(ctx context.Context, uid int64, biz string, bizID int64) error
	// DecrLike 移除点赞信息
	DecrLike(ctx context.Context, uid int64, biz string, bizID int64) error
	AddFavoriteItem(ctx context.Context, uid int64, favoriteID int64, biz string, bizID int64) error
	Get(ctx context.Context, uid int64, biz string, bizID int64) (domain.Interaction, error)

	// GetUserLikeInfo 获取用户的某个资源的点赞信息
	GetUserLikeInfo(ctx context.Context, uid int64, biz string, bizID int64) (domain.UserLike, error)
	// GetUserFavoriteInfo 获取用户的某个资源的收藏信息
	GetUserFavoriteInfo(ctx context.Context, uid int64, biz string, bizID int64) (domain.UserFavorite, error)
}

func (repo *interactionRepository) GetUserLikeInfo(ctx context.Context, uid int64, biz string, bizID int64) (domain.UserLike, error) {
	userLike, err := repo.dao.GetUserLikeInfo(ctx, uid, biz, bizID)
	if err != nil {
		return domain.UserLike{}, err
	}
	return domain.UserLike{
		ID:    userLike.ID,
		Biz:   userLike.Biz,
		Uid:   userLike.Uid,
		BizID: userLike.BizID,
		CTime: time.UnixMilli(userLike.CTime),
		UTime: time.UnixMilli(userLike.UTime),
	}, nil
}

func (repo *interactionRepository) GetUserFavoriteInfo(ctx context.Context, uid int64, biz string, bizID int64) (domain.UserFavorite, error) {
	entity, err := repo.dao.GetUserFavoriteInfo(ctx, uid, biz, bizID)
	if err != nil {
		return domain.UserFavorite{}, err
	}
	return domain.UserFavorite{
		ID:         entity.ID,
		Uid:        entity.Uid,
		Biz:        entity.Biz,
		BizID:      entity.BizID,
		FavoriteID: entity.FavoriteID,
		CTime:      time.UnixMilli(entity.CTime),
		UTime:      time.UnixMilli(entity.UTime),
	}, nil
}

func (repo *interactionRepository) Get(ctx context.Context, uid int64, biz string, bizID int64) (domain.Interaction, error) {
	// 因为是否点赞、是否收藏这俩个属性加入到Interaction领域对象中，所以从语义上应该从repository中来完成该对象的完整构造。
	inter, err := repo.cache.Get(ctx, biz, bizID)
	if err == nil {
		return repo.appendLikedAndCollected(ctx, inter, uid, biz, bizID), nil
	}

	interEntity, err := repo.dao.Get(ctx, biz, bizID)
	if err != nil {
		return inter, err
	}
	inter = repo.toDomain(interEntity)

	go func() {
		err := repo.cache.Set(ctx, biz, bizID, inter)
		if err != nil {
			// 记录日志，告警。
		}
	}()

	return repo.appendLikedAndCollected(ctx, inter, uid, biz, bizID), nil
}

func (repo *interactionRepository) appendLikedAndCollected(ctx context.Context,
	inter domain.Interaction, uid int64, biz string, bizID int64) domain.Interaction {
	var eg errgroup.Group
	eg.Go(func() error {
		userLike, err := repo.GetUserLikeInfo(ctx, uid, biz, bizID)
		if err != nil {
			// 可以考虑在这里记录日志
			inter.Liked = false
		} else {
			inter.Liked = userLike.Liked()
		}
		return nil
	})
	eg.Go(func() error {
		userFavorite, err := repo.GetUserFavoriteInfo(ctx, uid, biz, bizID)
		if err != nil {
			// 可以考虑在这里记录日志
			inter.Collected = false
		} else {
			inter.Collected = userFavorite.Collected()
		}
		return nil
	})
	_ = eg.Wait()
	return inter
}

func (repo *interactionRepository) AddFavoriteItem(ctx context.Context, uid int64, favoriteID int64, biz string, bizID int64) error {
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

func (repo *interactionRepository) BatchIncrReadCnt(ctx context.Context, bizs []string, bizIDs []int64) error {
	if len(bizIDs) != len(bizIDs) {
		return errors.New("the length of a and b must be equal")
	}
	// 1. 调用dao的批量增加
	err := repo.dao.BatchIncrReadCnt(ctx, bizs, bizIDs)
	if err != nil {
		return err
	}

	// 2. 批量增加缓存
	go func() {
		for i := range bizIDs {
			err := repo.cache.IncrReadCnt(ctx, bizs[i], bizIDs[i])
			if err != nil {
				// 记录日志
			}
		}
	}()
	return nil
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

func (repo *interactionRepository) toDomain(entity dao.Interaction) domain.Interaction {
	return domain.Interaction{
		ID:    entity.ID,
		Biz:   entity.Biz,
		BizID: entity.ID,

		Favorites: entity.Favorites,
		Views:     entity.ReadCnt,
		Likes:     entity.Likes,

		CTime: time.UnixMilli(entity.CTime),
		UTime: time.UnixMilli(entity.UTime),
	}
}
