package mongo

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"testing"
	"time"
)

// omitempty 如果字段存在omitempty，使用结构体构建filter以及进行更新时会忽略零值的字段。
type Article struct {
	Id       int64  `bson:"id,omitempty"`
	Title    string `bson:"title,omitempty"`
	Content  string `bson:"content,omitempty"`
	AuthorId int64  `bson:"author_id,omitempty"`
	Status   uint8  `bson:"status,omitempty"`
	Ctime    int64  `bson:"c_time,omitempty"`
	// 更新时间
	Utime int64 `bson:"u_time,omitempty"`
}

func TestCURD(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	monitor := event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Println(startedEvent.Command)
		},
	}
	opts := options.Client().ApplyURI("mongodb://root:example@localhost:27017").SetMonitor(&monitor)

	client, err := mongo.Connect(opts)
	assert.NoError(t, err)

	collection := client.Database("webook").Collection("articles")

	// 插入
	now := time.Now().UnixMilli()
	res, err := collection.InsertOne(ctx, Article{
		Id:      1,
		Title:   "my test",
		Content: "this is a test",
		Ctime:   now,
		Utime:   now,
	})
	assert.NoError(t, err)
	fmt.Printf("Inserted ID: %s \n", res.InsertedID)

	// 查询
	var art Article
	filter := bson.D{bson.E{Key: "id", Value: 1}}
	err = collection.FindOne(ctx, filter).Decode(&art)
	if err == mongo.ErrNoDocuments {
		fmt.Println("找不到id=1的article")
	} else {
		assert.NoError(t, err)
		fmt.Printf("article: %v\n", art)
	}

	// 更新
	sets := bson.D{bson.E{Key: "$set", Value: bson.E{Key: "title", Value: "new title."}}}
	updateRes, err := collection.UpdateMany(ctx, filter, sets)
	assert.NoError(t, err)
	fmt.Println("更新行数：", updateRes.ModifiedCount)

	// 删除
	defer func() {
		delRes, err := collection.DeleteMany(ctx, bson.D{})
		assert.NoError(t, err)
		fmt.Println("deleted count: ", delRes.DeletedCount)
	}()
}
