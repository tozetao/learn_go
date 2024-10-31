package ioc

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"learn_go/webook/internal/repository/dao"
)

func NewDB() *gorm.DB {
	type Config struct {
		DSN string
	}

	var config Config

	err := viper.UnmarshalKey("db", &config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("db config: %#v\n", config)

	db, err := gorm.Open(mysql.Open(config.DSN), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = dao.InitTable(db)
	if err != nil {
		panic("Init table failed.")
	}

	return db
}
