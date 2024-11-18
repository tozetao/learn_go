package dao

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Article struct {
	ID      int64  `gorm:"primaryKey,authIncrement"`
	Title   string `gorm:"type=varchar(1024)"`
	Content string `gorm:"type:blob"`
	Status  int8   `gorm:"type:tinyint"`

	AuthorID int64 `gorm:"index"`

	Ctime int64 `json:"c_time" gorm:"column:c_time"`
	Utime int64 `json:"u_time" gorm:"column:u_time"`
}

type PublishArticle struct {
	Article
}

type ArticleDao interface {
	Insert(ctx context.Context, data Article) (int64, error)
	UpdateByID(ctx context.Context, article Article) error
	Sync(ctx context.Context, article Article) (int64, error)
	SyncStatus(ctx context.Context, id int64, authorID int64, status int8) error
}

type GORMArticleDao struct {
	db *gorm.DB
}

func NewArticleDao(db *gorm.DB) ArticleDao {
	return &GORMArticleDao{
		db: db,
	}
}

func (dao *GORMArticleDao) Sync(ctx context.Context, article Article) (int64, error) {
	var (
		id  = article.ID
		err error
	)
	err = dao.db.Transaction(func(tx *gorm.DB) error {
		var err error
		authorDao := NewArticleAuthorDao(tx)
		if article.ID > 0 {
			err = authorDao.UpdateByID(ctx, article)
		} else {
			id, err = authorDao.Insert(ctx, article)
		}
		if err != nil {
			return err
		}
		article.ID = id

		pubArt := PublishArticle{Article: article}
		now := time.Now().UnixMilli()
		pubArt.Ctime = now
		pubArt.Utime = now

		err = tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"title":   pubArt.Title,
				"content": pubArt.Content,
				"u_time":  pubArt.Utime,
				"status":  pubArt.Status,
			}),
		}).Create(&pubArt).Error
		return err
	})
	return id, err
}

func (dao *GORMArticleDao) SyncStatus(ctx context.Context, id int64, authorID int64, status int8) error {
	now := time.Now().UnixMilli()

	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新制作库
		res := tx.Model(&Article{}).Where("id=? and author_id=?", id, authorID).
			Updates(map[string]interface{}{
				"status": status,
				"u_time": now,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			// 用户不可能执行到这里，需要记录日志
			return errors.New("failed to update article status")
		}

		// 更新线上库
		return tx.Model(&PublishArticle{}).Where("id=?", id).Updates(map[string]interface{}{
			"status": status,
			"u_time": now,
		}).Error
	})
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
	res := dao.db.WithContext(ctx).Model(&Article{}).
		Where("id = ? and author_id=?", article.ID, article.AuthorID).
		Updates(map[string]any{
			"title":   article.Title,
			"content": article.Content,
			"u_time":  article.Utime,
			"status":  article.Status,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("unauthorized operation, article_id: %v author_id: %v", article.ID, article.AuthorID)
	}
	return nil
}
