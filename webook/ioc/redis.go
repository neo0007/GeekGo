package ioc

import (
	"Neo/Workplace/goland/src/GeekGo/webook/config"
	"github.com/redis/go-redis/v9"
)

func InitRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr, // 对应 Docker Compose 中 Redis 的端口
		Password: "",                       // Redis 没有设置密码
		DB:       0,                        // 使用默认的 Redis DB
	})
}
