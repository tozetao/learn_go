package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"learn_go/webook/internal/domain"
)

type Interaction struct {
	ID int64 `json:"id" gorm:"primaryKey"`

	ArticleID int64 `json:"article_id"`
	ReadCnt   int   `json:"read_cnt"`
	Likes     int
	Favorites int

	CTime int64 `json:"c_time"`
	UTime int64 `json:"u_time"`
}

type InteractionDao interface {
	IncrReadCnt(ctx context.Context, inter domain.Interaction) error

	InsertLikeInfo(ctx context.Context, articleID int64) error
	DeleteLikeInfo(ctx context.Context, articleID int64) error
}

func (dao interactionDao) IncrReadCnt(ctx context.Context, inter domain.Interaction) error {
	return dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{}),
	}).Create(inter).Error
}

func (dao interactionDao) InsertLikeInfo(ctx context.Context, articleID int64) error {
	//TODO implement me
	panic("implement me")
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
