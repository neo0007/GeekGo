package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(paths string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, paths)
	return l
}
func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 不需要登录校验的
		for _, path := range l.paths {
			if c.Request.URL.Path == path {
				return
			}
		}

		sess := sessions.Default(c)
		userId := sess.Get("userId")
		if userId == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		sess.Set("userId", userId)
		sess.Options(sessions.Options{
			MaxAge: 30,
		})
		updateTime := sess.Get("updateTime")
		now := time.Now().UnixMilli()
		//如果没有刷新过
		if updateTime == nil {
			sess.Set("updateTime", now)
			sess.Save()
			return
		}
		//如果updateTime存在
		updateTimeVal, ok := updateTime.(int64)
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if now-updateTimeVal > 10*1000 {
			sess.Set("updateTime", now)
			sess.Save()
		}

	}
}
