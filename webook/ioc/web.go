package ioc

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/web"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/web/middleware"
	"Neo/Workplace/goland/src/GeekGo/webook/pkg/ginx/middleware/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitGin(mdls []gin.HandlerFunc, hdl *web.UserHandler) *gin.Engine {
	r := gin.Default()
	r.Use(mdls...)
	hdl.RegisterRoutes(r)
	return r
}

func InitMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		middleware.NewLoginJWTMiddlewareBuilder().
			IgnorePath("/users/login").
			IgnorePath("/users/login_sms/code/send").
			IgnorePath("/users/login_sms").
			IgnorePath("/users/signup").Build(),
		ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
	}

}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
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
	})
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
}
