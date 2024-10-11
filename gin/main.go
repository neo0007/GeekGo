package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello go!")
		//c.JSON(200, gin.H{
		//	"message": "pong",
		//})
	})

	r.POST("/post", func(c *gin.Context) {
		c.String(http.StatusOK, "hello post 方法")
	})

	r.GET("/users/:name", func(c *gin.Context) {
		name := c.Param("name")
		//c.JSON(http.StatusOK, gin.H{"hello, param router name": name})
		c.String(http.StatusOK, "hello, param router name: "+name)
	})

	r.GET("/views/*.a", func(c *gin.Context) {
		page := c.Param(".a")
		c.String(http.StatusOK, "hello, 这里是通配符路由: "+page)
	})

	r.GET("/order", func(c *gin.Context) {
		oid := c.Query("id")
		c.String(http.StatusOK, "hello, oid: "+oid)
	})

	r.Run("localhost:8080") // 监听并在 0.0.0.0:8080 上启动服务
}
