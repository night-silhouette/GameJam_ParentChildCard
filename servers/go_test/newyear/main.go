package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	// 增加一个简单的测试路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "服务器是好的"})
	})
	r.GET("/index", func(c *gin.Context) {
		c.File("./index.html")
	})
	r.Run(":8080")
}
