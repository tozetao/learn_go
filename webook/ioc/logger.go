package ioc

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"learn_go/webook/pkg/logger"
)

func NewLogger() logger.LoggerV2 {
	config := zap.NewDevelopmentConfig()
	err := viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
	l := zap.Must(config.Build())
	return logger.NewLogger(l)
}
