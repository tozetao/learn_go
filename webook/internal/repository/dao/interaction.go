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
	ReadCnt   int    `json:"read_cnt"`
	Likes     int
	Favorites int

	CTime int64 `json:"c_time"`
	UTime int64 `json:"u_time"`
}

const (
	LikeStatusUnknown = iota
	Liked
	Unliked
)

type ArticleLike struct {
	ID int64 `json:"id" gorm:"primaryKey,autoIncrement"`

	// 查询用户点赞的视频。
	// where uid = ? and biz = 'article'
	Uid   int64  `json:"uid" gorm:"index:idx_uid_biz_biz_id,unique"`
	Biz   string `json:"biz" gorm:"index:idx_uid_biz_biz_id,unique"`
	BizID int64  `json:"biz_id" gorm:"index:idx_uid_biz_biz_id,unique"`

	// 1表示点赞，0表示取消点赞
	Status uint8 `json:"status" gorm:"tinyint"`

	UTime int64 `json:"u_time" gorm:"column:u_time"`
	CTime int64 `json:"c_time" gorm:"column:c_time"`
}

type InteractionDao interface {
	IncrReadCnt(ctx context.Context, biz string, bizID int64) error

	InsertLikeInfo(ctx context.Context, articleID int64) error
	DeleteLikeInfo(ctx context.Context, articleID int64) error
}

func (dao interactionDao) IncrReadCnt(ctx context.Context, biz string, bizID int64) error {
	now := time.Now().UnixMilli()

	inter := Interaction{
		Biz:       biz,
		BizID:     bizID,
		ReadCnt:   1,
		Likes:     0,
		Favorites: 0,
		CTime:     now,
		UTime:     now,
	}

	return dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"u_time":   now,
			"read_cnt": gorm.Expr("read_cnt + ?", 1),
		}),
	}).Create(inter).Error
}

func (dao interactionDao) InsertLikeInfo(ctx context.Context, articleLike ArticleLike, inter Interaction) error {
	// 1. 插入点赞信息
	now := time.Now().UnixMilli()
	articleLike.CTime = now
	articleLike.UTime = now
	articleLike.Status = Liked
	
	err := dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"u_time":   now,
			""
		}),
	}).Create(&articleLike).Error
}

func (dao interactionDao) DeleteLikeInfo(ctx context.Context, articleID int64) error {
	//TODO implement me
	panic("implement me")
}

type interactionDao struct {
	db *gorm.DB
}

func NewInteractionDao(db *gorm.DB) InteractionDao {
	return &interactionDao{
		db: db,
	}
}
