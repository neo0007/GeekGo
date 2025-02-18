package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

type jwtHandler struct {
	//access_token key
	atKey []byte
	//refresh_token key
	rtKey []byte
}

func newJwtHandler() jwtHandler {
	return jwtHandler{
		atKey: []byte("56j6wp8hlc8biryjns2ju2n6g02f6fyu"),
		rtKey: []byte("56j6wp8hlc8biryjns2ju2n6g02f6ldu"),
	}
}

func (h jwtHandler) setJWTToken(c *gin.Context, uid int64) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 120)),
		},
		Uid:       uid,
		UserAgent: c.Request.UserAgent(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(h.atKey)
	if err != nil {
		c.String(http.StatusInternalServerError, "系统错误")
		return err
	}
	c.Header("x-jwt-token", tokenStr)
	return nil
}

func (h jwtHandler) setRefreshToken(c *gin.Context, uid int64) error {
	claims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1200)),
		},
		Uid: uid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(h.rtKey)
	if err != nil {
		c.String(http.StatusInternalServerError, "系统错误")
		return err
	}
	c.Header("x-refresh-token", tokenStr)
	return nil
}

func ExtractToken(c *gin.Context) string {
	tokenHeader := c.GetHeader("Authorization")
	segs := strings.Split(tokenHeader, " ")
	if len(segs) != 2 {
		//没登录，有人瞎搞
		c.AbortWithStatus(http.StatusUnauthorized)
		return ""
	}
	return segs[1]
}

type RefreshClaims struct {
	Uid int64
	jwt.RegisteredClaims
}

type UserClaims struct {
	jwt.RegisteredClaims
	//声明你要放进token里面的数据
	Uid       int64
	UserAgent string
}
