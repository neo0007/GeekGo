package ioc

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms/failover"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms/memory"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms/ratelimit"
	pkgRatelimit "Neo/Workplace/goland/src/GeekGo/webook/pkg/ratelimit"
	"github.com/redis/go-redis/v9"
	"time"
)

//func InitSMSService() sms.Service {
//	return memory.NewService()
//}

//func InitSMSService(cmdable redis.Cmdable) sms.Service {
//	// 创建原始短信服务
//	baseService := memory.NewService()
//
//	// 创建限流器
//	limiter := pkgRatelimit.NewRedisSlidingWindowLimiter(cmdable, time.Second, 100)
//
//	// 用限流装饰器包装原始短信服务
//	return ratelimit.NewRatelimitSMSService(baseService, limiter)
//}

func InitSMSService(cmdable redis.Cmdable) sms.Service {
	// 1. 创建原始短信服务
	svc1 := memory.NewService() // 第一个短信服务
	svc2 := memory.NewService() // 第二个短信服务（可以换成不同的实现，比如其他提供商）

	// 2. 创建限流器
	limiter1 := pkgRatelimit.NewRedisSlidingWindowLimiter(cmdable, time.Second, 100)
	limiter2 := pkgRatelimit.NewRedisSlidingWindowLimiter(cmdable, time.Second, 100)

	// 3. 用限流装饰器包装短信服务
	rateLimitedSvc1 := ratelimit.NewRatelimitSMSService(svc1, limiter1)
	rateLimitedSvc2 := ratelimit.NewRatelimitSMSService(svc2, limiter2)

	// 4. 组合失败转移策略
	return failover.NewFailoverSMSService([]sms.Service{rateLimitedSvc1, rateLimitedSvc2})
}
