package main

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository/dao"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/web"
	"github.com/gin-contrib/cors"
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
	r.Use(func(c *gin.Context) {
		println("this is second middleware")
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

	initDB := false
	if initDB {
		err = dao.InitTable(db)
		if err != nil {
			panic(err)
		}
	}

	return db
}
