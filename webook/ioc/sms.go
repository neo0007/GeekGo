package ioc

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms/auth"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms/failover"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms/logger"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms/memory"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms/ratelimit"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms/retryable"
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

//func InitSMSService(cmdable redis.Cmdable) sms.Service {
//	// 1. 创建原始短信服务
//	svc1 := memory.NewService() // 第一个短信服务
//	svc2 := memory.NewService() // 第二个短信服务（可以换成不同的实现，比如其他提供商）
//
//	// 2. 创建限流器
//	limiter1 := pkgRatelimit.NewRedisSlidingWindowLimiter(cmdable, time.Second, 100)
//	limiter2 := pkgRatelimit.NewRedisSlidingWindowLimiter(cmdable, time.Second, 100)
//
//	// 3. 用限流装饰器包装短信服务
//	rateLimitedSvc1 := ratelimit.NewRatelimitSMSService(svc1, limiter1)
//	rateLimitedSvc2 := ratelimit.NewRatelimitSMSService(svc2, limiter2)
//
//	// 4. 组合失败转移策略
//	failoverService := failover.NewFailoverSMSService([]sms.Service{rateLimitedSvc1, rateLimitedSvc2})
//
//	// 5. 认证装饰器
//	secretKey := "my-secret-key" // 业务方的密钥
//	authService := auth.NewAuthSMSService(failoverService, secretKey)
//
//	return authService
//}

func InitSMSService(cmdable redis.Cmdable) sms.Service {
	// 1. 创建原始短信服务
	svc1 := memory.NewService() // 第一个短信服务
	svc2 := memory.NewService() // 第二个短信服务（可以换成不同的实现）

	// 2. 创建限流器
	limiter1 := pkgRatelimit.NewRedisSlidingWindowLimiter(cmdable, time.Second, 100)
	limiter2 := pkgRatelimit.NewRedisSlidingWindowLimiter(cmdable, time.Second, 100)

	// 3. 先加限流装饰器
	rateLimitedSvc1 := ratelimit.NewRatelimitSMSService(svc1, limiter1)
	rateLimitedSvc2 := ratelimit.NewRatelimitSMSService(svc2, limiter2)

	// 4. 失败转移策略（failover）
	failoverService := failover.NewFailoverSMSService([]sms.Service{rateLimitedSvc1, rateLimitedSvc2})

	//5. 加入重试机制
	retryMax := 3 // 最大重试次数
	retryService := retryable.NewService(failoverService, retryMax)

	//6. 认证装饰器
	secretKey := []byte("56j6wp8hlc8biryjns2ju2n6g02f6fyu") // 业务方的密钥
	authService := auth.NewAuthSMSService(retryService, secretKey)

	//// 7. 同步转异步机制
	//asyncService := async.NewSMSService(authService, repo)
	//
	//// 启动异步 goroutine 处理未发送的请求
	//asyncService.StartAsync()

	//return asyncService

	// 7. **日志装饰器**
	logService := logger.NewService(authService)

	return logService
}

func InitSMSServiceV1(cmdable redis.Cmdable) sms.Service {
	return memory.NewService()
}
