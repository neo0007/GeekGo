//go:build k8s

package config

var Config = config{
	DB: DBConfig{
		// k8s 连接
		DSN: "root:root@tcp(webook-mysql:3308)/webook",
	},
	Redis: RedisConfig{
		Addr: "webook-redis:6380",
	},
}
