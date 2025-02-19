package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler interface {
	ExtractToken(c *gin.Context) string
	SetLoginToken(c *gin.Context, uid int64) error
	SetJWTToken(c *gin.Context, uid int64, ssid string) error
	CheckSession(c *gin.Context, ssid string) error
	ClearToken(c *gin.Context) error
}

type RefreshClaims struct {
	Uid  int64
	Ssid string
	jwt.RegisteredClaims
}

type UserClaims struct {
	jwt.RegisteredClaims
	//声明你要放进token里面的数据
	Uid       int64
	Ssid      string
	UserAgent string
}
