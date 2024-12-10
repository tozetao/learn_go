package viper_test

import (
	"fmt"
	"github.com/spf13/viper"
	"testing"
)

func TestViper(t *testing.T) {
	type Metric struct {
		Host string
		Port int
	}
	type Warehouse struct {
		Host1 string `mapstructure:"host"`
		Port1 int    `mapstructure:"port"`
	}
	type DataStore struct {
		Metric    Metric
		Warehouse Warehouse
	}

	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../webook/config")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	var config DataStore
	err := viper.UnmarshalKey("datastore", &config)
	if err != nil {
		panic(err)
	}

	fmt.Printf("metric: %v, warehouse: %v\n", config.Metric, config.Warehouse)
}

func TestViper2(t *testing.T) {
	type Config struct {
		Addrs []string
	}
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../webook/config")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	var config Config
	err := viper.UnmarshalKey("kafka", &config)
	if err != nil {
		panic(err)
	}
	fmt.Println(config)
}
