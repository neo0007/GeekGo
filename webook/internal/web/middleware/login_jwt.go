package middleware

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePath(paths string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, paths)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 不需要登录校验的
		for _, path := range l.paths {
			if c.Request.URL.Path == path {
				return
			}
		}

		tokenHeader := c.GetHeader("Authorization")
		if tokenHeader == "" {
			//没登录
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			//没登录，有人瞎搞
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		claims := &web.UserClaims{}
		// ParseWithClaims 参数一定要传指针，因为它会修改参数内容，如果是结构体本身，它只会复制一份新的，你拿不到修改后的数据

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("56j6wp8hlc8biryjns2ju2n6g02f6fyu"), nil
		})
		if err != nil {
			//没登录
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//err 为 nil，token 不为 nil
		if !token.Valid || claims.Uid == 0 {
			//没登录
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		now := time.Now()
		// 每 10秒刷新一次
		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString([]byte("56j6wp8hlc8biryjns2ju2n6g02f6fyu"))
			if err != nil {
				//记录日志
				log.Println("jwt 续约失败", err)
			}
			c.Header("x-jwt-token", tokenStr)
		}

		c.Set("claims", claims)
	}
}
