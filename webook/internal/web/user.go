package web

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/domain"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/service"
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

// UserHandler 定义与用户有关的路由
type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern    = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
		passwordRegexPattern = "^[A-Za-z\\d@$!%*?&]{8,}$"
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
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
	if errors.Is(err, service.ErrUserDuplicateEmail) {
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
	if err := c.ShouldBindJSON(&req); err != nil {
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
	sess.Options(sessions.Options{
		MaxAge: 30,
	})
	sess.Save()

	c.String(http.StatusOK, "登录成功！")
	return

}

func (u *UserHandler) LoginJWT(c *gin.Context) {
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
		c.String(http.StatusOK, "用户名或密码不对！")
		return
	}
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		Uid:       user.Id,
		UserAgent: c.Request.UserAgent(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("56j6wp8hlc8biryjns2ju2n6g02f6fyu"))
	if err != nil {
		c.String(http.StatusInternalServerError, "系统错误")
		return
	}
	c.Header("x-jwt-token", tokenStr)

	fmt.Println(tokenStr)
	fmt.Println(user)

	c.String(http.StatusOK, "登录成功！")
	return

}

func (u *UserHandler) Edit(c *gin.Context) {

}

func (u *UserHandler) Profile(c *gin.Context) {
	c.String(http.StatusOK, "这是你的 Profile！")
}

func (u *UserHandler) ProfileJWT(c *gin.Context) {
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
	println(claims.Uid)
	c.String(http.StatusOK, "这是你的 Profile")
}

type UserClaims struct {
	jwt.RegisteredClaims
	//声明你要放进token里面的数据
	Uid       int64
	UserAgent string
}
