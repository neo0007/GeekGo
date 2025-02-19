package jwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strings"
	"time"
)

var (
	AtKey = []byte("56j6wp8hlc8biryjns2ju2n6g02f6fyu")
	RtKey = []byte("56j6wp8hlc8biryjns2ju2n6g02f6ldu")
)

type RedisJWTHandler struct {
	cmd redis.Cmdable
}

func NewRedisJWTHandler(cmd redis.Cmdable) Handler {
	return &RedisJWTHandler{
		cmd: cmd,
	}
}

func (h RedisJWTHandler) ExtractToken(c *gin.Context) string {
	tokenHeader := c.GetHeader("Authorization")
	segs := strings.Split(tokenHeader, " ")
	if len(segs) != 2 {
		//没登录，有人瞎搞
		c.AbortWithStatus(http.StatusUnauthorized)
		return ""
	}
	return segs[1]
}

func (h RedisJWTHandler) SetLoginToken(c *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.SetJWTToken(c, uid, ssid)
	if err != nil {
		return err
	}
	err = h.setRefreshToken(c, uid, ssid)
	return err
}

func (h RedisJWTHandler) setRefreshToken(c *gin.Context, uid int64, ssid string) error {
	claims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		Uid:  uid,
		Ssid: ssid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(RtKey)
	if err != nil {
		c.String(http.StatusInternalServerError, "系统错误")
		return err
	}
	c.Header("x-refresh-token", tokenStr)
	return nil
}

func (h RedisJWTHandler) CheckSession(c *gin.Context, ssid string) error {
	_, err := h.cmd.Exists(c, fmt.Sprintf("user:ssid:%s", ssid)).Result()
	return err
}

func (h RedisJWTHandler) ClearToken(c *gin.Context) error {
	c.Header("x-jwt-token", "")
	c.Header("x-refresh-token", "")
	cs, ok := c.Get("claims")
	if !ok {
		//可以考虑监控这里
		c.String(http.StatusOK, "系统错误")
		return errors.New("claims not found")
	}
	// ok代表是不是 *UserClaims
	claims, ok := cs.(*UserClaims)
	if !ok {
		//可以考虑监控这里
		c.String(http.StatusOK, "系统错误")
		return errors.New("claims error")
	}
	return h.cmd.Set(c, fmt.Sprintf("user:ssid:%s", claims.Ssid),
		"", time.Hour*24*7).Err()
}

func (h RedisJWTHandler) SetJWTToken(c *gin.Context, uid int64, ssid string) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 120)),
		},
		Uid:       uid,
		Ssid:      ssid,
		UserAgent: c.Request.UserAgent(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(AtKey)
	if err != nil {
		c.String(http.StatusInternalServerError, "系统错误")
		return err
	}
	c.Header("x-jwt-token", tokenStr)
	return nil
}
