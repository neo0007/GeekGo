package ioc

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms/memory"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms/ratelimit"
	pkgRatelimit "Neo/Workplace/goland/src/GeekGo/webook/pkg/ratelimit"
	"github.com/redis/go-redis/v9"
	"time"
)

//func InitSMSService() sms.Service {
//	return memory.NewService()
//}

func InitSMSService(cmdable redis.Cmdable) sms.Service {
	// 创建原始短信服务
	baseService := memory.NewService()

	// 创建限流器
	limiter := pkgRatelimit.NewRedisSlidingWindowLimiter(cmdable, time.Second, 100)

	// 用限流装饰器包装原始短信服务
	return ratelimit.NewRatelimitSMSService(baseService, limiter)
}
