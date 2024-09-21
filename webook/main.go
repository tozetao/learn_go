package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	templateId := "test-template"
	server := InitWebServer(templateId)

	server.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome!")
	})
	server.Run(":8100")
}

//func initServer() *gin.Engine {
//	server := gin.Default()
//
//	//// rate limiter
//	//rdb := initRedis()
//	//server.Use(ratelimit.NewBuilder(rdb, time.Second*10, 100).Build())
//
//	//// 设置session中间件
//	//store, err := redisstore.NewStore(6, "tcp", "192.168.1.100:6379", "",
//	//	[]byte("8yF7u3sG4hJkZbQeRtDpNxVmCiLwOa9H"), []byte("qWeR7tYvUiOpKjLzHaBxNcDmFsAg4R5E"))
//	//if err != nil {
//	//	panic("init redis store failed.")
//	//}
//	//server.Use(sessions.Sessions("ssid", store))
//
//	// 用户认证（session实现）
//	//login := middleware.NewLoginMiddlewareBuilder()
//	//s := []string{"/users/login", "/users/signup", "/test1", "/test/1", "/test2"}
//	//server.Use(login.IgnorePath(s...).Builder())
//
//	return server
//}
