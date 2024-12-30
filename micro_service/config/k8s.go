//go:build k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(webook-mysql:4399)/webook",
	},
	Redis: RedisConfig{
		Addr: "webook-redis:14001",
	},
}
