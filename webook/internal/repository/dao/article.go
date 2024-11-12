package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Article struct {
	ID      int64  `gorm:"primaryKey,authIncrement"`
	Title   string `gorm:"type=varchar(1024)"`
	Content string `gorm:"type:blob"`

	AuthorID int64 `gorm:"index"`

	Ctime int64 `json:"c_time" gorm:"column:c_time"`
	Utime int64 `json:"u_time" gorm:"column:u_time"`
}

type ArticleDao interface {
	Insert(ctx context.Context, data Article) (int64, error)
	UpdateByID(ctx context.Context, article Article) error
}

type GORMArticleDao struct {
	db *gorm.DB
}

func NewArticleDao(db *gorm.DB) ArticleDao {
	return &GORMArticleDao{
		db: db,
	}
}

func (dao *GORMArticleDao) Insert(ctx context.Context, article Article) (int64, error) {
	now := time.Now()

	article.Ctime = now.UnixMilli()
	article.Utime = now.UnixMilli()

	err := dao.db.WithContext(ctx).Create(&article).Error
	return article.ID, err
}

func (dao *GORMArticleDao) UpdateByID(ctx context.Context, article Article) error {
	now := time.Now().UnixMilli()
	article.Utime = now

	// 依赖gorm，忽略零值的更新，会用主键进行更新。
	// 问题：可读性很差。
	// dao.db.WithContext(ctx).Updates(&article)

	// tip:
	// 通过ID更新帖子。一般都是更新帖子的内容，id和作者id肯定是对应的，因此方法可以命名为UpdateByID，不要UpdateByIDAndAuthorID
	res := dao.db.WithContext(ctx).Model(&article).
		Where("id = ? and author_id=?", article.ID, article.AuthorID).
		Updates(map[string]any{
			"title":   article.Title,
			"content": article.Content,
			"u_time":  article.Utime,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("unauthorized operation, article_id: %v author_id: %v", article.ID, article.AuthorID)
	}
	return nil
}
