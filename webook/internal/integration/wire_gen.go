//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject

package integration

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

// Injectors from wire.go:

func InitWebServer(templateId string) *gin.Engine {
	v := ioc.InitMiddlewares()
	smsService := ioc.InitSMSService()
	cmdable := ioc.NewRedis()
	codeCache := cache.NewCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	codeService := service.NewCodeService(templateId, smsService, codeRepository)
	smsHandler := web.NewSMSHandler(codeService)
	db := ioc.NewDB()
	userDao := dao.NewUserDao(db)
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDao, userCache)
	userService := service.NewUserService(userRepository)
	userHandler := web.NewUserHandler(userService, codeService)
	engine := ioc.InitGin(v, smsHandler, userHandler)
	return engine
}

// wire.go:

var (
	providers = wire.NewSet(ioc.NewDB, ioc.NewRedis, cache.NewCodeCache, cache.NewUserCache, dao.NewUserDao, repository.NewCodeRepository, repository.NewUserRepository, ioc.InitSMSService, service.NewCodeService, service.NewUserService, web.NewSMSHandler, web.NewUserHandler, ioc.InitMiddlewares, ioc.InitGin)
)
