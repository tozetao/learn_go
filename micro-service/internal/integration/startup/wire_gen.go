// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	article2 "learn_go/webook/internal/event/article"
	"learn_go/webook/internal/job"
	"learn_go/webook/internal/repository"
	"learn_go/webook/internal/repository/article"
	"learn_go/webook/internal/repository/cache"
	"learn_go/webook/internal/repository/dao"
	"learn_go/webook/internal/service"
	"learn_go/webook/internal/web"
	"learn_go/webook/ioc"
)

// Injectors from wire.go:

func InitArticleHandler() *web.ArticleHandler {
	db := NewDB()
	articleDao := dao.NewArticleDao(db)
	cmdable := NewRedis()
	articleCache := cache.NewArticleCache(cmdable)
	userDao := dao.NewUserDao(db)
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDao, userCache)
	loggerV2 := ioc.NewLogger()
	articleRepository := article.NewArticleRepository(articleDao, articleCache, userRepository, loggerV2)
	authorRepository := article.NewArticleAuthorRepository()
	readerRepository := article.NewArticleReaderRepository()
	config := ioc.NewSaramaConfig()
	syncProducer := ioc.NewSyncProducer(config)
	producer := article2.NewSyncProducer(syncProducer)
	articleService := service.NewArticleService(articleRepository, authorRepository, readerRepository, producer, loggerV2)
	interactionDao := dao.NewInteractionDao(db)
	interactionCache := cache.NewInteractionCache(cmdable)
	interactionRepository := repository.NewInteractionRepository(interactionDao, interactionCache)
	interactionService := service.NewInteractionService(interactionRepository)
	articleHandler := web.NewArticleHandler(articleService, interactionService, loggerV2)
	return articleHandler
}

func InitArticleHandlerV1(articleDao dao.ArticleDao) *web.ArticleHandler {
	cmdable := NewRedis()
	articleCache := cache.NewArticleCache(cmdable)
	db := NewDB()
	userDao := dao.NewUserDao(db)
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDao, userCache)
	loggerV2 := ioc.NewLogger()
	articleRepository := article.NewArticleRepository(articleDao, articleCache, userRepository, loggerV2)
	authorRepository := article.NewArticleAuthorRepository()
	readerRepository := article.NewArticleReaderRepository()
	config := ioc.NewSaramaConfig()
	syncProducer := ioc.NewSyncProducer(config)
	producer := article2.NewSyncProducer(syncProducer)
	articleService := service.NewArticleService(articleRepository, authorRepository, readerRepository, producer, loggerV2)
	interactionDao := dao.NewInteractionDao(db)
	interactionCache := cache.NewInteractionCache(cmdable)
	interactionRepository := repository.NewInteractionRepository(interactionDao, interactionCache)
	interactionService := service.NewInteractionService(interactionRepository)
	articleHandler := web.NewArticleHandler(articleService, interactionService, loggerV2)
	return articleHandler
}

func InitWebServer(templateId string) *gin.Engine {
	cmdable := NewRedis()
	jwtHandler := web.NewJWTHandler(cmdable)
	loggerV2 := ioc.NewLogger()
	v := ioc.InitMiddlewares(jwtHandler, loggerV2)
	smsService := ioc.InitSMSService()
	codeCache := cache.NewCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	codeService := service.NewCodeService(templateId, smsService, codeRepository)
	smsHandler := web.NewSMSHandler(codeService)
	db := NewDB()
	userDao := dao.NewUserDao(db)
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDao, userCache)
	userService := service.NewUserService(userRepository, loggerV2)
	userHandler := web.NewUserHandler(userService, codeService, jwtHandler)
	oAuth2Service := ioc.InitOAuth2Service()
	oAuth2WechatHandler := web.NewOAuth2WechatHandler(oAuth2Service, userService, jwtHandler)
	engine := ioc.InitGin(v, smsHandler, userHandler, oAuth2WechatHandler)
	return engine
}

func InitScheduler() *job.Scheduler {
	db := NewDB()
	jobDao := dao.NewJobDao(db)
	jobRepository := repository.NewCronJobRepository(jobDao)
	loggerV2 := ioc.NewLogger()
	jobService := service.NewJobService(jobRepository, loggerV2)
	scheduler := job.NewScheduler(jobService, loggerV2)
	return scheduler
}

// wire.go:

var (
	providers = wire.NewSet(ioc.NewLogger, NewDB,
		NewRedis, cache.NewCodeCache, cache.NewUserCache, dao.NewUserDao, repository.NewCodeRepository, repository.NewUserRepository, ioc.InitSMSService, ioc.InitOAuth2Service, service.NewCodeService, service.NewUserService, web.NewSMSHandler, web.NewUserHandler, web.NewOAuth2WechatHandler, web.NewJWTHandler, web.NewTestHandler, ioc.InitMiddlewares, ioc.InitGin,
	)

	articleProviders = wire.NewSet(
		NewDB,
		NewRedis, ioc.NewLogger, ioc.NewSaramaConfig, ioc.NewConsumerClient, article2.NewBatchReadEventConsumer, ioc.NewConsumers, ioc.NewSyncProducer, article2.NewSyncProducer, web.NewArticleHandler, service.NewArticleService, service.NewInteractionService, repository.NewInteractionRepository, dao.NewInteractionDao, cache.NewInteractionCache, article.NewArticleRepository, dao.NewArticleDao, cache.NewArticleCache, repository.NewUserRepository, dao.NewUserDao, cache.NewUserCache, article.NewArticleAuthorRepository, article.NewArticleReaderRepository,
	)

	articleProvidersV1 = wire.NewSet(
		NewDB,
		NewRedis, ioc.NewLogger, ioc.NewSaramaConfig, ioc.NewConsumerClient, article2.NewBatchReadEventConsumer, ioc.NewConsumers, ioc.NewSyncProducer, article2.NewSyncProducer, web.NewArticleHandler, service.NewArticleService, service.NewInteractionService, repository.NewInteractionRepository, dao.NewInteractionDao, cache.NewInteractionCache, article.NewArticleRepository, cache.NewArticleCache, repository.NewUserRepository, dao.NewUserDao, cache.NewUserCache, article.NewArticleAuthorRepository, article.NewArticleReaderRepository,
	)

	schedulerProvider = wire.NewSet(
		NewDB,
		NewRedis, ioc.NewLogger, job.NewScheduler, service.NewJobService, repository.NewCronJobRepository, dao.NewJobDao,
	)
)