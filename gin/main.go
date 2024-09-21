package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/items/*.html", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ok")
	})
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
