package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	v1 := "2024-10-01 11:00:10"
	layout := "2006-01-02 15:04:05"

	t1, err := time.ParseInLocation(layout, v1, time.Local)
	assert.NoError(t, err)

	fmt.Printf("%v\n", t1)
}

func TestGormFind(t *testing.T) {
	db := NewDB()

	// 查询文章列表
	var arts []Article
	res := db.Model(&Article{}).Where("author_id = ?", 1000).Offset(0).Limit(10).Find(&arts)
	fmt.Printf("%v, %v\n", res.Error, res.RowsAffected)
	fmt.Printf("%v\n", arts)
}

func NewDB() *gorm.DB {
	dsn := "root:root@tcp(127.0.0.1:3306)/webook"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{LogLevel: logger.Info, SlowThreshold: time.Second}),
	})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&Article{})
	if err != nil {
		panic(err)
	}

	return db
}

type Article struct {
	ID      int64  `gorm:"primaryKey,authIncrement" bson:"id,omitempty"`
	Title   string `gorm:"type=varchar(1024)"  bson:"title,omitempty"`
	Content string `gorm:"type:blob"  bson:"content,omitempty"`
	Status  int8   `gorm:"type:tinyint"  bson:"status,omitempty"`

	AuthorID int64 `gorm:"index"  bson:"author_id,omitempty"`

	Ctime int64 `json:"c_time" gorm:"column:c_time"  bson:"c_time,omitempty"`
	Utime int64 `json:"u_time" gorm:"column:u_time"  bson:"u_time,omitempty"`
}

//func InitTable(db *gorm.DB) error {
//	return db.AutoMigrate(&dao.User{}, &Article{}, &PublishArticle{})
//}
