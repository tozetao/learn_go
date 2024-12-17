package ioc

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
	"learn_go/webook/internal/repository/dao"
	"learn_go/webook/pkg/logger"
)

func NewDB(log logger.LoggerV2) *gorm.DB {
	type Config struct {
		DSN string
	}

	var config Config

	err := viper.UnmarshalKey("db", &config)
	if err != nil {
		panic(err)
	}
	log.Info("db config", logger.Field{Key: "", Value: config})

	db, err := gorm.Open(mysql.Open(config.DSN), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = dao.InitTable(db)
	if err != nil {
		panic("Init table failed.")
	}

	// 接入Prometheus
	err = db.Use(prometheus.New(prometheus.Config{
		DBName:          "webook",
		RefreshInterval: 15,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{"thread_running"},
			},
		},
	}))
	if err != nil {
		panic(err)
	}

	return db
}
