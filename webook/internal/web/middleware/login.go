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
		updateTime := sess.Get("updateTime")
		if updateTime == nil || userId == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 当执行sessions.save 时，需要重新 set所有要保存的值，原来的数据将会覆盖所以这里要重新设置一下userId
		sess.Set("userId", userId)
		//session.Options需要 Save 方法后才能发挥作用，第一次进入的时候没有 Save，
		//需要在login函数中调用控制一次，这里是刷新时候用到的，重新 Save 后如果不设置
		//将不在保留原来Save 的内容
		sess.Options(sessions.Options{
			MaxAge: 30,
		})
		now := time.Now().UnixMilli()
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
