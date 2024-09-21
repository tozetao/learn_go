package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"learn_go/webook/internal/web"
	"learn_go/webook/internal/web/middleware"
	"strings"
	"time"
)

func InitGin(middlewares []gin.HandlerFunc,
	smsHandler *web.SMSHandler, userHandler *web.UserHandler,
	oauthWechatHandler *web.OAuth2WechatHandler) *gin.Engine {
	server := gin.Default()
	server.Use(middlewares...)

	smsHandler.RegisterRoutes(server)
	userHandler.RegisterRoutes(server)
	oauthWechatHandler.RegisterRoutes(server)
	return server
}

func InitMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		authHdl(),
	}
}

// CORSMiddleware CORS中间件
func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		// 允许跨域的源
		// AllowOrigins:     []string{"https://foo.com"},
		// 允许跨域的方法
		// AllowMethods:     []string{"PUT", "PATCH"},
		// 跨域时允许携带的请求头
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		// 跨域时允许读取的响应头
		ExposeHeaders: []string{"x-jwt-token"},
		// 是否允许携带cookie
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "x.com")
		},
		MaxAge: 12 * time.Hour,
	})
}

// authHdl 用户验证（jwt）
func authHdl() gin.HandlerFunc {
	login := middleware.NewLoginJWTMiddlewareBuilder()
	s := []string{
		"/users/login",
		"/users/signup",
		"/",
		"/demo",
		"/users/login_sms",
		"/sms/send",
	}
	return login.IgnorePath(s...).Builder()
}

//func initLimiter() gin.HandlerFunc {
//	l := limiter.NewRedisSideWindow()
//	builder := ratelimit.NewBuilder("limit:ip:")
//}
