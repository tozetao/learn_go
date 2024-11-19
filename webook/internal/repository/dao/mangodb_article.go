package dao

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// MangoDBArticleDao Article的mongo存储实现
type MangoDBArticleDao struct {
	db              *mongo.Database
	artCol          *mongo.Collection
	publishedArtCol *mongo.Collection
	node            *snowflake.Node
}

func NewMongoArticleDao(db *mongo.Database) ArticleDao {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	return &MangoDBArticleDao{
		db:              db,
		artCol:          db.Collection("articles"),
		publishedArtCol: db.Collection("published_articles"),
		node:            node,
	}
}

func (dao *MangoDBArticleDao) Insert(ctx context.Context, article Article) (int64, error) {
	article.ID = dao.node.Generate().Int64()
	now := time.Now().UnixMilli()
	article.Ctime = now
	article.Utime = now

	// 先插入制作库
	_, err := dao.artCol.InsertOne(ctx, article)
	if err != nil {
		return 0, err
	}
	// 再插入线上库
	pubArt := PublishArticle{
		Article: article,
	}
	_, err = dao.publishedArtCol.InsertOne(ctx, pubArt)
	if err != nil {
		return 0, err
	}
	return article.ID, nil
}

func (dao *MangoDBArticleDao) UpdateByID(ctx context.Context, article Article) error {
	// 先更新制作库
	now := time.Now().UnixMilli()
	article.Utime = now
	filter := bson.D{{"id", article.ID}, {"author_id", article.AuthorID}}
	set := bson.D{
		{"$set", bson.D{
			{"title", article.Title},
			{"content", article.Content},
			{"u_time", now},
		}},
	}
	_, err := dao.artCol.UpdateOne(ctx, filter, set)
	if err != nil {
		return err
	}
	_, err = dao.publishedArtCol.UpdateOne(ctx, bson.D{
		{"id", article.ID},
	}, set)
	if err != nil {
		return err
	}
	return nil
}

func (dao *MangoDBArticleDao) Sync(ctx context.Context, article Article) (int64, error) {
	panic("implement me")
}

func (dao *MangoDBArticleDao) SyncStatus(ctx context.Context, id int64, authorID int64, status int8) error {
	panic("implement me")
}
