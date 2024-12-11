//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	event "learn_go/webook/internal/event/article"
	"learn_go/webook/internal/repository"
	"learn_go/webook/internal/repository/article"
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

		ioc.InitMiddlewares,
		ioc.InitGin,

		ioc.InitSMSService,
		ioc.InitOAuth2Service,

		// 消费者
		ioc.NewSaramaConfig,
		ioc.NewConsumerClient,
		//event.NewConsumer,
		event.NewBatchReadEventConsumer,
		ioc.NewConsumers,
		// 生产者
		ioc.NewSyncProducer,
		event.NewSyncProducer,

		web.NewSMSHandler,
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		web.NewJWTHandler,
		web.NewArticleHandler,
		web.NewTestHandler,

		service.NewCodeService,
		service.NewUserService,
		service.NewArticleService,
		service.NewInteractionService,

		repository.NewInteractionRepository,
		repository.NewCodeRepository,
		repository.NewUserRepository,
		article.NewArticleRepository,
		article.NewArticleReaderRepository,
		article.NewArticleAuthorRepository,

		dao.NewUserDao,
		dao.NewInteractionDao,
		dao.NewArticleDao,

		cache.NewArticleCache,
		cache.NewCodeCache,
		cache.NewUserCache,
		cache.NewInteractionCache,

		wire.Struct(new(App), "*"),
	)
)

func InitApp(templateId string) *App {
	wire.Build(providers)
	return new(App)
}

//func InitWebServer(templateId string) *gin.Engine {
//	wire.Build(providers)
//	return gin.Default()
//}
