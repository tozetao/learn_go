//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"learn_go/webook/internal/repository"
	"learn_go/webook/internal/repository/cache"
	"learn_go/webook/internal/repository/dao"
	"learn_go/webook/internal/service"
	"learn_go/webook/internal/web"
	"learn_go/webook/ioc"
)

var (
	providers = wire.NewSet(
		// 第三方依赖
		ioc.NewLogger,
		ioc.NewDB,
		ioc.NewRedis,

		cache.NewCodeCache, cache.NewUserCache,

		dao.NewUserDao,

		repository.NewCodeRepository,
		repository.NewUserRepository,

		ioc.InitSMSService,
		ioc.InitOAuth2Service,
		service.NewCodeService,
		service.NewUserService,

		web.NewSMSHandler,
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		web.NewJWTHandler,
		web.NewTestHandler,

		ioc.InitMiddlewares,
		ioc.InitGin,
	)
)

func InitWebServer(templateId string) *gin.Engine {
	wire.Build(providers)
	return gin.Default()
}
