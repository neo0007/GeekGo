package web

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/oauth2/wechat"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
	"time"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	jwtHandler
	stateKey []byte
	cfg      WechatHandlerConfig
}

type WechatHandlerConfig struct {
	Secure bool
}

func NewOAuth2WechatHandler(svc wechat.Service, userSvc service.UserService, cfg WechatHandlerConfig) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:        svc,
		userSvc:    userSvc,
		stateKey:   []byte("56j6wp8hlc8biryjns2ju2n6g02f9fyu"),
		cfg:        cfg,
		jwtHandler: newJwtHandler(),
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", h.AuthURL)
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	state := uuid.New()
	url, err := h.svc.AuthURL(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "构造扫码登录 URL 失败",
		})
		return
	}

	err = h.setStateCookie(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})
}

func (h *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		State: state,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Minute)),
		},
	})
	tokenStr, err := token.SignedString(h.stateKey)
	if err != nil {
		return fmt.Errorf("token.SignedString: %w", err)
	}
	ctx.SetCookie("jwt-state", tokenStr, 600,
		"/oauth2/wechat/callback", "", h.cfg.Secure, true)
	return nil
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	//验证微信的 code
	code := ctx.Query("code")
	err := h.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "登录失败",
		})
		return
	}

	info, err := h.svc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	// 从 userSvc中拿到 uid
	u, err := h.userSvc.FindOrCreateByWechat(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	err = h.setJWTToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	err = h.setRefreshToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

func (h *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	ck, err := ctx.Cookie("jwt-state")
	if err != nil {
		//有人搞你，
		//做好监控
		//记录日志
		return fmt.Errorf("拿不到state的cookie: %w", err)
	}
	var sc StateClaims
	token, err := jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return h.stateKey, nil
	})
	if err != nil || token.Valid != true {
		//做好监控
		//记录日志
		return fmt.Errorf("token解析错误或过期: %w", err)
	}
	if state != sc.State {
		//做好监控
		//记录日志
		return errors.New("State不相等")
	}
	return nil
}

type StateClaims struct {
	State string
	jwt.RegisteredClaims
}
