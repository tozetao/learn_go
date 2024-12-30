package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type ArticleAuthorDao interface {
	UpdateByID(ctx context.Context, article Article) error

	Insert(ctx context.Context, article Article) (int64, error)
}

func (dao *articleAuthorDao) UpdateByID(ctx context.Context, article Article) error {
	now := time.Now().UnixMilli()
	res := dao.db.WithContext(ctx).Model(&article).
		Where("id = ? and author_id = ?", article.ID, article.AuthorID).
		Updates(map[string]interface{}{
			"title":   article.Title,
			"content": article.Content,
			"status":  article.Status,
			"u_time":  now,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		// 要记录日志
		return errors.New("failed to update, no affected rows")
	}
	return nil
}

func (dao *articleAuthorDao) Insert(ctx context.Context, article Article) (int64, error) {
	now := time.Now().UnixMilli()
	article.Ctime = now
	article.Utime = now
	err := dao.db.WithContext(ctx).Create(&article).Error
	return article.ID, err
}

type articleAuthorDao struct {
	db *gorm.DB
}

func NewArticleAuthorDao(db *gorm.DB) ArticleAuthorDao {
	return &articleAuthorDao{
		db: db,
	}
}
