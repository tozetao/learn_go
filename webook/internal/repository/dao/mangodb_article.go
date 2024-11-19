package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

// MangoDBArticleDao Article的mongo存储实现
type MangoDBArticleDao struct {
	db  *mongo.Database
	col *mongo.Collection
}

func NewMongoArticleDao(db *mongo.Database) ArticleDao {
	return &MangoDBArticleDao{
		db:  db,
		col: db.Collection("articles"),
	}
}

func (dao *MangoDBArticleDao) Insert(ctx context.Context, data Article) (int64, error) {
	panic("implement me")
}

func (dao *MangoDBArticleDao) UpdateByID(ctx context.Context, article Article) error {
	panic("implement me")
}

func (dao *MangoDBArticleDao) Sync(ctx context.Context, article Article) (int64, error) {
	panic("implement me")
}

func (dao *MangoDBArticleDao) SyncStatus(ctx context.Context, id int64, authorID int64, status int8) error {
	panic("implement me")
}
