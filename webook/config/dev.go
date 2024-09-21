//go:build !k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(127.0.0.1:3306)/webook",
	},
	Redis: RedisConfig{
		Addr: "192.168.1.100:6379",
	},
}
