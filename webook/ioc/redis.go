package ioc

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	addr := viper.GetString("redis.addr")
	//println(addr)
	return redis.NewClient(&redis.Options{
		Addr:     addr, // 对应 Docker Compose 中 Redis 的端口
		Password: "",   // Redis 没有设置密码
		DB:       0,    // 使用默认的 Redis DB
	})
}
