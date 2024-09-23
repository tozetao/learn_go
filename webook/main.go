package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func main() {
	//configFile := pflag.String("config", "config/dev.yaml", "配置文件路径")
	//pflag.Parse()
	//fmt.Printf("configFile: %s\n", *configFile)

	//TestViper()

	server := InitWebServer("test-template")
	server.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "hello world")
	})
	server.Run(":9130")
}

func InitConfigV1() {
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	type Config struct {
		Addr string
	}
	c := Config{}
	err = viper.UnmarshalKey("redis", &c)
	if err != nil {
		panic(fmt.Errorf("unmarshal error: %w", err))
	}

	type DBConfig struct {
		DSN string
	}
	var dbConfig DBConfig
	err = viper.UnmarshalKey("db", &dbConfig)
	fmt.Printf("%v\n", dbConfig)
}

// 测试下etcd
//
