package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ArticleReaderDao interface {
	Upsert(ctx context.Context, article Article) error
}

type articleReaderDao struct {
	db *gorm.DB
}

func (dao *articleReaderDao) Upsert(ctx context.Context, article Article) error {
	// Clauses用于添加额外的SQL子句到查询中。
	// clause.OnConflict: 指定发生冲突时的行为。在这里是当列id发生冲突时，应该更新title和content。
	return dao.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   article.Title,
			"content": article.Content,
		}),
	}).Create(&article).Error
}

func NewArticleReaderDao() ArticleReaderDao {
	return &articleReaderDao{}
}
