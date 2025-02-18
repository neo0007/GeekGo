//go:build wireinject

package main

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository/cache/redis"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository/dao/gorm"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/web"
	"Neo/Workplace/goland/src/GeekGo/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitDB, ioc.InitRedis,
		// 初始化 DAO
		gorm.NewUserDao,

		redis.NewUserCache,
		redis.NewCodeCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,

		service.NewUserService,
		service.NewCodeService,
		// 基于内存的实现
		ioc.InitSMSService,
		ioc.InitWechatService,
		ioc.NewWechatHandlerConfig,

		web.NewUserHandler,
		web.NewOAuth2WechatHandler,

		ioc.InitGin,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}
