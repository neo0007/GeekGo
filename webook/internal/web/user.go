package web

import "github.com/gin-gonic/gin"

// UserHandler 定义与用户有关的路由
type UserHandler struct {
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
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.Profile)
}

func (u *UserHandler) Signup(c *gin.Context) {

}

func (u *UserHandler) Login(c *gin.Context) {

}

func (u *UserHandler) Edit(c *gin.Context) {

}

func (u *UserHandler) Profile(c *gin.Context) {

}
