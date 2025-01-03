package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func main() {
	loadConfig()
	//InitLogger()

	app := InitApp("test-template")

	// 启动定时任务
	app.cron.Start()

	// 启动监控服务
	initPrometheus()

	// 启动web服务
	app.server.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "hello world")
	})
	app.server.Run(":9130")

	// 等web服务器关闭后，再关系其他服务
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	select {
	case <-ctx.Done():
	// 有任务长时间执行，超时关闭
	case <-app.cron.Stop().Done():
		//	job都执行结束，正常关闭
	}
}

func loadConfig() {
	//configFile := pflag.String("config", "./config/dev.yaml", "配置文件路径")
	//pflag.Parse()
	//fmt.Printf("%v\n", *configFile)
	//viper.SetConfigFile(*configFile)

	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
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

func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8081", nil)
	}()
}
