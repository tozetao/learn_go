package ioc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"learn_go/webook/config"
	"learn_go/webook/internal/repository/dao"
)

func NewDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = dao.InitTable(db)
	if err != nil {
		panic("Init table failed.")
	}

	return db
}
