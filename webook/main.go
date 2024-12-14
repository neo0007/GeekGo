package main

import (
	"Neo/Workplace/goland/src/GeekGo/webook/config"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository/cache"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository/dao"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms/memory"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/web"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/web/middleware"
	"Neo/Workplace/goland/src/GeekGo/webook/pkg/ginx/middleware/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	db := initDB()
	rdb := initRedis()
	u := initUser(db, rdb)

	r := initWebServer()

	//注册路由须在中间件 cors 跨域插件运行之后，否则不会生效！
	u.RegisterRoutes(r)

	//r := gin.Default()
	//
	//r.GET("/hello", func(c *gin.Context) {
	//	c.String(http.StatusOK, "hello world")
	//})

	r.Run(":8081")
	//err := r.Run("localhost:8080")
	//if err != nil {
	//	panic(err)
	//}
}

func initWebServer() *gin.Engine {
	r := gin.Default()

	// r.Use 的作用是全局使用，即应用在全部路由上生效，首先执行的是 Use，然后执行路由
	r.Use(func(c *gin.Context) {
		println("this is first middleware")
	})

	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	r.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	r.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"http://localhost:3000"},
		//AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		//如果省略上面 AllowMethods:...... 则所有 POST、GET等全部方法将都被允许
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 你不加ExposeHeaders，前端是拿不到对应数据的
		ExposeHeaders:    []string{"Content-Length", "Authorization", "x-jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "webook.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	// cookie 的设置也需要放在解决跨域问题之后，否则也会出现不可预料的错误！
	//store := cookie.NewStore([]byte("secret"))
	// 单机用 memstore:
	//store := memstore.NewStore([]byte("56j6wp8hlc8biryjns2ju2n6g02f6fyu"),
	//[]byte("jp74g2x60gqqv2mrn36xpzmussrmyeyx"))
	//多实例部署用 redis:
	//store, err := redis.NewStore(16,
	//	"tcp", "localhost:6379", "",
	//	[]byte("56j6wp8hlc8biryjns2ju2n6g02f6fyu"),
	//	[]byte("jp74g2x60gqqv2mrn36xpzmussrmyeyx"))
	//if err != nil {
	//	panic(err)
	//}
	//r.Use(sessions.Sessions("mysession", store))

	r.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePath("/users/login").
		IgnorePath("/users/login_sms/code/send").
		IgnorePath("/users/login_sms").
		IgnorePath("/users/signup").Build())

	return r
}

func initUser(db *gorm.DB, rdb redis.Cmdable) *web.UserHandler {
	ud := dao.NewUserDao(db)
	uc := cache.NewUserCache(rdb)
	repo := repository.NewUserRepository(ud, uc)
	svc := service.NewUserService(repo)
	codeCache := cache.NewCodeCache(rdb)
	codeRepo := repository.NewCodeRepository(codeCache)
	smsSvc := memory.NewService()
	codeSvc := service.NewCodeService(codeRepo, smsSvc)
	u := web.NewUserHandler(svc, codeSvc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		// 应该只在初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		// 一旦初始化出错，即退出
		panic(err)
	}

	initDB := true
	if initDB {
		err = dao.InitTable(db)
		if err != nil {
			panic(err)
		}
	}

	return db
}

func initRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr, // 对应 Docker Compose 中 Redis 的端口
		Password: "",                       // Redis 没有设置密码
		DB:       0,                        // 使用默认的 Redis DB
	})
}

//func (*LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		//不需要校验
//		if c.Request.URL.Path == "/users/login" ||
//			c.Request.URL.Path == "/users/signup" {
//			return
//		}
//
//	}

//}
