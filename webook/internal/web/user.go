package web

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/domain"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service"
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

const biz = "login"

// 确保 handler 上实现了 UserHandler 的接口
var _ handler = (*UserHandler)(nil)

// UserHandler 定义与用户有关的路由
type UserHandler struct {
	svc         service.UserService
	codeSvc     service.CodeService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService) *UserHandler {
	const (
		emailRegexPattern    = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
		passwordRegexPattern = "^[A-Za-z\\d@$!%*?&]{8,}$"
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		codeSvc:     codeSvc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}

}

//func (u *UserHandler) RegisterRoutes(r *gin.Engine) {
//	r.POST("/users/signup", u.Signup)
//	r.POST("/users/login", u.Login)
//	r.POST("/users/edit", u.Edit)
//	r.GET("/users/profile", u.Profile)
//}

func (u *UserHandler) RegisterRoutes(r *gin.Engine) {
	ug := r.Group("/users")
	ug.POST("/signup", u.Signup)
	//ug.POST("/login", u.Login)
	ug.POST("login", u.LoginJWT)
	ug.POST("/edit", u.Edit)
	//ug.GET("/profile", u.Profile)
	ug.GET("/profile", u.ProfileJWT)
	ug.POST("/login_sms/code/send", u.SendSMSLoginCode)
	ug.POST("/login_sms", u.LoginSMS)
}

func (u *UserHandler) SendSMSLoginCode(c *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Phone == "" {
		c.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "输入有误",
		})
		return
	}
	var err = u.codeSvc.Send(c, biz, req.Phone)
	switch err {
	case nil:
		c.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		c.JSON(http.StatusOK, Result{
			Msg: "发送验证码太频繁, 1分钟后再试",
		})
	default:
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
}

func (u *UserHandler) LoginSMS(c *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, Result{
			Code: 5,
			Msg:  "bind error",
		})
		return
	}
	ok, err := u.codeSvc.Verify(c, biz, req.Phone, req.Code)
	if errors.Is(err, service.ErrCodeVerifyTooManyTimes) {
		c.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证次数太多，请重新发送验证码",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		c.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码有误",
		})
		return
	}

	user, err := u.svc.FindOrCreate(c, req.Phone)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "DB系统错误",
		})
		return
	}

	if err = u.setJWTToken(c, user.Id); err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	c.JSON(http.StatusOK, Result{
		Code: 5,
		Msg:  "验证码校验通过",
	})
}

func (u *UserHandler) Signup(c *gin.Context) {
	type SignupReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	//c.String(http.StatusOK, "signup successfully before")
	var req SignupReq
	//Bind 方法会根据 Content-Type 来解析你的数据到 req 里面
	//解析错了，就会直接写会一个 400 的错误
	//if err := c.Bind(&req); err != nil {
	//	return
	//}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		c.String(http.StatusOK, "系统错误！")
		return
	}
	if !ok {
		c.String(http.StatusOK, "你的邮箱格式不对！")
		return
	}

	if req.ConfirmPassword != req.Password {
		c.String(http.StatusOK, "两次输入的密码不一致！")
		return
	}

	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		//记录日志
		c.String(http.StatusOK, "系统错误！")
		return
	}
	if !ok {
		c.String(http.StatusOK, "密码必须大于 8 位")
		return
	}

	//调用一下 svc 的方法, 存储数据到数据库
	err = u.svc.Signup(c, domain.User{Email: req.Email, Password: req.Password})
	if errors.Is(err, service.ErrUserDuplicate) {
		c.String(http.StatusOK, err.Error())
		return
	}
	if err != nil {
		c.String(http.StatusOK, "系统异常")
		return
	}

	c.String(http.StatusOK, "signup successfully")
}

func (u *UserHandler) Login(c *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := u.svc.Login(c, domain.User{Email: req.Email, Password: req.Password})
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		c.String(http.StatusOK, "用户名或密码不对！")
		return
	}
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	sess := sessions.Default(c)
	sess.Set("userId", user.Id)
	sess.Set("updateTime", time.Now().UnixMilli())
	sess.Save()

	c.String(http.StatusOK, "登录成功！")
}

func (u *UserHandler) LoginJWT(c *gin.Context) {
	//这一句防止原来已设置的 jwt token 仍然在有效期，从而出现已登录的错误, 下面这句不管用，后面再检查测试
	//c.Header("Authorization", "Bearer")

	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := u.svc.Login(c, domain.User{Email: req.Email, Password: req.Password})
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		//下面为重置“x-jwt-token”, 避免原有token在有效期内直接进入profile
		//请注意如果设置"x-jwt-token"值为空，即“”，则不发生作用，请设置成任何非空字符串即可
		c.Header("x-jwt-token", "reset")
		c.String(http.StatusOK, "用户名或密码不对！")
		return
	}
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	if err = u.setJWTToken(c, user.Id); err != nil {
		c.String(http.StatusOK, "系统错误！")
		return
	}

	c.String(http.StatusOK, "登录成功！")
}

func (u *UserHandler) setJWTToken(c *gin.Context, uid int64) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 120)),
		},
		Uid:       uid,
		UserAgent: c.Request.UserAgent(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("56j6wp8hlc8biryjns2ju2n6g02f6fyu"))
	if err != nil {
		c.String(http.StatusInternalServerError, "系统错误")
		return err
	}
	c.Header("x-jwt-token", tokenStr)
	return nil
}

func (u *UserHandler) Edit(c *gin.Context) {

}

func (u *UserHandler) Profile(c *gin.Context) {
	c.String(http.StatusOK, "这是你的 Profile！")
}

func (u *UserHandler) ProfileJWT(c *gin.Context) {
	type userProfile struct {
		Email    string `json:"Email"`
		Phone    string `json:"Phone"`
		Nickname string `json:"Nickname"`
		Birthday string `json:"Birthday"`
		AboutMe  string `json:"AboutMe"`
	}

	cs, ok := c.Get("claims")
	if !ok {
		//可以考虑监控这里
		c.String(http.StatusOK, "系统错误")
		return
	}
	// ok代表是不是 *UserClaims
	claims, ok := cs.(*UserClaims)
	if !ok {
		//可以考虑监控这里
		c.String(http.StatusOK, "系统错误")
		return
	}
	user, err := u.svc.Profile(c, claims.Uid)
	if err != nil {
		c.String(http.StatusOK, err.Error())
	}

	var uProfile = userProfile{
		Email:    user.Email,
		Phone:    user.Phone,
		Nickname: user.Nickname,
		Birthday: user.Birthday,
		AboutMe:  user.AboutMe,
	}

	c.JSON(http.StatusOK, uProfile)
	//println(claims.Uid)
	//c.String(http.StatusOK, "这是你的 Profile")
}

type UserClaims struct {
	jwt.RegisteredClaims
	//声明你要放进token里面的数据
	Uid       int64
	UserAgent string
}
