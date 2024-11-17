package startup

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"learn_go/webook/internal/repository/dao"
)

func NewDB() *gorm.DB {
	dsn := "root:root@tcp(127.0.0.1:3306)/webook"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{LogLevel: logger.Info, SlowThreshold: time.Second}),
	})
	if err != nil {
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}

	return db
}
