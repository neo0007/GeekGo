package main

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository/dao"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/web"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	db := initDB()
	u := initUser(db)

	r := initWebServer()

	// 注册路由须在中间件 cors 跨域插件运行之后，否则不会生效！
	u.RegisterRoutes(r)

	err := r.Run("localhost:8080")
	if err != nil {
		panic(err)
	}
}

func initWebServer() *gin.Engine {
	r := gin.Default()

	// r.Use 的作用是全局使用，即应用在全部路由上生效，首先执行的是 Use，然后执行路由
	r.Use(func(c *gin.Context) {
		println("this is first middleware")
	})

	r.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"http://localhost:3000"},
		//AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		//如果省略上面 AllowMethods:...... 则所有 POST、GET等全部方法将都被允许
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	// cookie 的设置也需要放在解决跨域问题之后，否则也会出现不可预料的错误！
	//store := cookie.NewStore([]byte("secret"))
	// 单机用 memstore, 多实例部署用 redis
	//store := memstore.NewStore([]byte("56j6wp8hlc8biryjns2ju2n6g02f6fyu"), []byte("jp74g2x60gqqv2mrn36xpzmussrmyeyx"))
	store, err := redis.NewStore(16,
		"tcp", "localhost:6379", "",
		[]byte("56j6wp8hlc8biryjns2ju2n6g02f6fyu"),
		[]byte("jp74g2x60gqqv2mrn36xpzmussrmyeyx"))
	if err != nil {
		panic(err)
	}
	r.Use(sessions.Sessions("mysession", store))

	r.Use(middleware.NewLoginMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login").Build())

	return r
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDao(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
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
