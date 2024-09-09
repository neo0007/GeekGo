package main

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/web"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	u := &web.UserHandler{}
	u.RegisterRoutes(r)

	r.Run("localhost:8080")
}
