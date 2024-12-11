package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Interaction struct {
	ID int64 `json:"id" gorm:"primaryKey,autoIncrement"`

	// biz, biz_id组成唯一索引
	Biz       string `json:"biz" gorm:"index:idx_biz_biz_id,unique"`
	BizID     int64  `json:"biz_id" gorm:"index:idx_biz_biz_id,unique"`
	ReadCnt   int64  `json:"read_cnt"`
	Likes     int64
	Favorites int64

	CTime int64 `json:"c_time"`
	UTime int64 `json:"u_time"`
}

const (
	LikeStatusUnknown = iota
	Liked
	Unliked
)

type UserLike struct {
	ID int64 `json:"id" gorm:"primaryKey,autoIncrement"`

	// 查询用户点赞的视频。
	// where uid = ? and biz = 'article'
	Uid   int64  `json:"uid" gorm:"index:idx_uid_biz_biz_id,unique"`
	Biz   string `json:"biz" gorm:"index:idx_uid_biz_biz_id,unique"`
	BizID int64  `json:"biz_id" gorm:"index:idx_uid_biz_biz_id,unique"`

	// 1表示点赞，2表示取消点赞
	Status uint8 `json:"status" gorm:"tinyint"`

	UTime int64 `json:"u_time" gorm:"column:u_time"`
	CTime int64 `json:"c_time" gorm:"column:c_time"`
}

type UserFavorite struct {
	ID int64 `json:"id" gorm:"primaryKey,autoIncrement"`

	Uid int64 `json:"uid" gorm:"uniqueIndex:idx_uid_biz_type_id"`

	Biz   string `json:"biz" gorm:"type:varchar(128);uniqueIndex:idx_uid_biz_type_id"`
	BizID int64  `json:"biz_id" gorm:"uniqueIndex:idx_uid_biz_type_id"`

	// 收藏夹id是唯一的，自己有索引
	FavoriteID int64 `json:"favorite_id" gorm:"index"`

	UTime int64 `json:"u_time" gorm:"column:u_time"`
	CTime int64 `json:"c_time" gorm:"column:c_time"`
}

type InteractionDao interface {
	IncrReadCnt(ctx context.Context, biz string, bizID int64) error

	InsertLikeInfo(ctx context.Context, uid int64, biz string, bizID int64) error
	DeleteLikeInfo(ctx context.Context, uid int64, biz string, bizID int64) error
	InsertFavorite(ctx context.Context, favorite UserFavorite) error
	Get(ctx context.Context, biz string, bizID int64) (Interaction, error)
	GetUserLikeInfo(ctx context.Context, uid int64, biz string, bizID int64) (UserLike, error)
	GetUserFavoriteInfo(ctx context.Context, uid int64, biz string, id int64) (UserFavorite, error)
	BatchIncrReadCnt(ctx context.Context, bizs []string, ds []int64) error
}

func (dao *interactionDao) GetUserFavoriteInfo(ctx context.Context, uid int64, biz string, bizID int64) (UserFavorite, error) {
	var userFavorite UserFavorite
	err := dao.db.WithContext(ctx).Model(&UserFavorite{}).
		Where("uid = ? and biz = ? and biz_id = ?", uid, biz, bizID).
		First(&userFavorite).Error
	return userFavorite, err
}

func (dao *interactionDao) GetUserLikeInfo(ctx context.Context, uid int64, biz string, bizID int64) (UserLike, error) {
	var userLike UserLike
	err := dao.db.WithContext(ctx).Model(&UserLike{}).
		Where("uid = ? and biz = ? and biz_id = ?", uid, biz, bizID).First(&userLike).Error
	return userLike, err
}

func (dao *interactionDao) Get(ctx context.Context, biz string, bizID int64) (Interaction, error) {
	var inter Interaction
	err := dao.db.WithContext(ctx).Model(&Interaction{}).
		Where("biz = ? and biz_id = ?", biz, bizID).First(&inter).Error
	return inter, err
}

func (dao *interactionDao) InsertFavorite(ctx context.Context, favorite UserFavorite) error {
	// 暂时不考虑分享数的正确性
	now := time.Now().UnixMilli()
	return dao.db.Transaction(func(tx *gorm.DB) error {
		// 新建收藏记录
		err := tx.WithContext(ctx).Create(&favorite).Error
		if err != nil {
			return err
		}
		//增加分享次数
		return tx.WithContext(ctx).Model(&Interaction{}).
			Clauses(clause.OnConflict{
				DoUpdates: clause.Assignments(map[string]interface{}{
					"favorites": gorm.Expr("favorites + 1"),
					"u_time":    now,
				}),
			}).Create(&Interaction{
			Biz:       favorite.Biz,
			BizID:     favorite.BizID,
			Favorites: 1,
			CTime:     now,
			UTime:     now,
		}).Error
	})
}

func (dao *interactionDao) IncrReadCnt(ctx context.Context, biz string, bizID int64) error {
	now := time.Now().UnixMilli()

	inter := Interaction{
		Biz:     biz,
		BizID:   bizID,
		ReadCnt: 1,
		CTime:   now,
		UTime:   now,
	}

	return dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"u_time":   now,
			"read_cnt": gorm.Expr("read_cnt + ?", 1),
		}),
	}).Create(&inter).Error
}

// InsertLikeInfo 插入点赞记录
func (dao *interactionDao) InsertLikeInfo(ctx context.Context, uid int64, biz string, bizID int64) error {
	/*
		目前这种实现，没有判定玩家是否点赞成功，因此只要用户不停的取消并再次点赞，点赞数就会一直增加。
		解决方案：
			锁住一条记录
			不存在则插入
			更新则更新
		并发高的时候会导致死锁，但是锁住的是个人数据，几率会低，可以接受。
	*/

	// 1. 插入点赞信息
	now := time.Now().UnixMilli()

	err := dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"u_time": now,
			// 因为使用status来表示记录是否存在，所以由dao层来设置status的值
			"status": Liked,
		}),
	}).Create(&UserLike{
		Uid:    uid,
		Biz:    biz,
		BizID:  bizID,
		Status: Liked,
		CTime:  now,
		UTime:  now,
	}).Error
	if err != nil {
		return err
	}

	// 增加计数
	return dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"u_time": now,
			"likes":  gorm.Expr("likes + ?", 1),
		}),
	}).Create(&Interaction{
		Biz:   biz,
		BizID: bizID,
		Likes: 1,
		CTime: now,
		UTime: now,
	}).Error
}

func (dao *interactionDao) DeleteLikeInfo(ctx context.Context, uid int64, biz string, bizID int64) error {
	now := time.Now().UnixMilli()

	return dao.db.Transaction(func(tx *gorm.DB) error {
		err := tx.WithContext(ctx).Model(&UserLike{}).
			Where("uid = ? and biz = ? and biz_id = ?", uid, biz, bizID).Updates(map[string]interface{}{
			"u_time": now,
			"status": Unliked,
		}).Error
		if err != nil {
			return err
		}
		// 严格一点的需要判定执行行数是否等于1，成功在能够取消点赞
		return tx.WithContext(ctx).Model(&Interaction{}).
			Where("biz = ? and biz_id = ?", biz, bizID).
			Updates(map[string]interface{}{
				"u_time": now,
				"likes":  gorm.Expr("likes - ?", 1),
			}).Error
	})
}

type interactionDao struct {
	db *gorm.DB
}

func (dao *interactionDao) BatchIncrReadCnt(ctx context.Context, bizs []string, ids []int64) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		for i := range ids {
			err := dao.IncrReadCnt(ctx, bizs[i], ids[i])
			if err != nil {
				// 对于文章这种计数器，即使少了个别文章的计数影响也不大，因此不回滚，只记录该错误。
			}
		}
		return nil
	})
}

func NewInteractionDao(db *gorm.DB) InteractionDao {
	return &interactionDao{
		db: db,
	}
}
