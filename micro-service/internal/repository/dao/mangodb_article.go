package dao

import (
	"context"
	"errors"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// MangoDBArticleDao Article的mongo存储实现
type MangoDBArticleDao struct {
	db              *mongo.Database
	artCol          *mongo.Collection
	publishedArtCol *mongo.Collection
	node            *snowflake.Node
}

func (dao *MangoDBArticleDao) GetByID(ctx context.Context, id int64) (Article, error) {
	//TODO implement me
	panic("implement me")
}

func (dao *MangoDBArticleDao) ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]Article, error) {
	//TODO implement me
	panic("implement me")
}

func (dao *MangoDBArticleDao) GetPubByID(ctx context.Context, id int64) (PublishArticle, error) {
	//TODO implement me
	panic("implement me")
}

func NewMongoArticleDao(db *mongo.Database, node *snowflake.Node) ArticleDao {
	return &MangoDBArticleDao{
		db:              db,
		artCol:          db.Collection("articles"),
		publishedArtCol: db.Collection("published_articles"),
		node:            node,
	}
}

func (dao *MangoDBArticleDao) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]Article, error) {
	panic("implement me")
}

func (dao *MangoDBArticleDao) Insert(ctx context.Context, article Article) (int64, error) {
	article.ID = dao.node.Generate().Int64()
	now := time.Now().UnixMilli()
	article.Ctime = now
	article.Utime = now

	_, err := dao.artCol.InsertOne(ctx, article)
	if err != nil {
		return 0, err
	}
	return article.ID, nil
}

func (dao *MangoDBArticleDao) UpdateByID(ctx context.Context, article Article) error {
	now := time.Now().UnixMilli()
	article.Utime = now
	filter := bson.D{{"id", article.ID}, {"author_id", article.AuthorID}}
	set := bson.D{
		{"$set", bson.D{
			{"title", article.Title},
			{"content", article.Content},
			{"status", article.Status},
			{"u_time", now},
		}},
	}
	res, err := dao.artCol.UpdateOne(ctx, filter, set)
	if err != nil {
		return err
	}
	if res.ModifiedCount != 1 {
		return errors.New("failed to update article")
	}
	return nil
}

func (dao *MangoDBArticleDao) Sync(ctx context.Context, article Article) (int64, error) {
	// 先更新制作库
	var (
		err error
		id  = article.ID
	)
	if article.ID > 0 {
		err = dao.UpdateByID(ctx, article)
	} else {
		id, err = dao.Insert(ctx, article)
	}
	if err != nil {
		return 0, err
	}

	article.ID = id

	now := time.Now().UnixMilli()
	pubArt := PublishArticle(article)
	pubArt.Ctime = now
	pubArt.Utime = now
	filter := bson.M{"id": id}
	set := bson.D{
		{"$setOnInsert", bson.D{
			{"id", pubArt.ID},
			{"author_id", pubArt.AuthorID},
			{"c_time", pubArt.Ctime},
		}},
		{"$set", bson.D{
			{"title", pubArt.Title},
			{"content", pubArt.Content},
			{"status", pubArt.Status},
			{"u_time", pubArt.Utime},
		}},
	}
	_, err = dao.publishedArtCol.UpdateOne(
		ctx, filter, set, options.Update().SetUpsert(true))
	if err != nil {
		return 0, err
	}
	return article.ID, nil
}

func (dao *MangoDBArticleDao) SyncStatus(ctx context.Context, id int64, authorID int64, status int8) error {
	// 先更新制作库
	var (
		err error
	)
	now := time.Now().UnixMilli()
	filter := bson.D{
		{"id", id},
		{"author_id", authorID},
	}
	set := bson.D{
		{"$set", bson.D{
			{"status", status},
			{"u_time", now},
		}},
	}
	artRes, err := dao.artCol.UpdateOne(ctx, filter, status)
	if err != nil {
		return err
	}
	if artRes.ModifiedCount != 1 {
		// id、author_id不一致，要记录错误日志，告警。
		return errors.New("failed to update status")
	}

	_, err = dao.publishedArtCol.UpdateOne(ctx, bson.M{
		"id": id,
	}, set)
	return err
}
