package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	InitConfig()
	//InitLogger()

	server := InitWebServer("test-template")
	server.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "hello world")
	})
	server.Run(":9130")
}

func InitConfig() {
	configFile := pflag.String("config", "config/dev.yaml", "配置文件路径")
	pflag.Parse()

	//viper.SetConfigName("dev")
	//viper.SetConfigType("yaml")
	//viper.AddConfigPath("config")

	viper.SetConfigFile(*configFile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func InitRemoteConfig() {
	err := viper.AddRemoteProvider("etcd3", "http://127.0.0.1:12379", "/webook")
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}

	str := viper.GetString("user1")
	fmt.Printf("user1: %s\n", str)
}

func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}
