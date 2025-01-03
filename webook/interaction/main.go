package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

/*
App

	组装grpc、kafka消费者等各种服务

配置文件的加载
kafka消费者服务
grpc服务

*/

func LoadConfig() {
	configFile := pflag.String("config", "./config/test.yaml", "配置文件路径")
	pflag.Parse()
	fmt.Printf("%v\n", *configFile)
	viper.SetConfigFile(*configFile)

	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func main() {
	LoadConfig()

	app := InitApp()

	// 启动消费者服务
	for _, consumer := range app.consumers {
		consumer.Start()
	}

	err := app.server.Start()
	log.Printf("err: %v\n", err)
}
