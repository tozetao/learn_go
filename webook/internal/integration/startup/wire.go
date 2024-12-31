//go:build wireinject
// +build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	event2 "learn_go/webook/interaction/event/article"
	repository2 "learn_go/webook/interaction/repository"
	cache2 "learn_go/webook/interaction/repository/cache"
	dao2 "learn_go/webook/interaction/repository/dao"
	service2 "learn_go/webook/interaction/service"
	event "learn_go/webook/internal/event/article"
	"learn_go/webook/internal/job"
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
		NewDB,
		NewRedis,

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

	articleProviders = wire.NewSet(
		NewDB,
		NewRedis,
		ioc.NewLogger,

		// 消费者
		ioc.NewSaramaConfig,
		ioc.NewConsumerClient,
		//event.NewConsumer,
		event2.NewBatchReadEventConsumer,
		ioc.NewConsumers,
		// 生产者
		ioc.NewSyncProducer,
		event.NewSyncProducer,

		web.NewArticleHandler,
		service.NewArticleService,
		service2.NewInteractionService,

		repository2.NewInteractionRepository,
		dao2.NewInteractionDao,
		cache2.NewInteractionCache,

		article.NewArticleRepository,

		dao.NewArticleDao,
		cache.NewArticleCache,
		repository.NewUserRepository,

		dao.NewUserDao,
		cache.NewUserCache,

		article.NewArticleAuthorRepository,
		article.NewArticleReaderRepository,
	)

	articleProvidersV1 = wire.NewSet(
		NewDB,
		NewRedis,
		ioc.NewLogger,

		// 消费者
		ioc.NewSaramaConfig,
		ioc.NewConsumerClient,
		//event.NewConsumer,
		event2.NewBatchReadEventConsumer,
		ioc.NewConsumers,
		// 生产者
		ioc.NewSyncProducer,
		event.NewSyncProducer,

		web.NewArticleHandler,
		service.NewArticleService,
		service2.NewInteractionService,

		repository2.NewInteractionRepository,
		dao2.NewInteractionDao,
		cache2.NewInteractionCache,

		article.NewArticleRepository,

		cache.NewArticleCache,
		repository.NewUserRepository,

		dao.NewUserDao,
		cache.NewUserCache,

		article.NewArticleAuthorRepository,
		article.NewArticleReaderRepository,
	)

	schedulerProvider = wire.NewSet(
		NewDB,
		NewRedis,
		ioc.NewLogger,

		job.NewScheduler,
		service.NewJobService,
		repository.NewCronJobRepository,
		dao.NewJobDao,
	)
)

func InitArticleHandler() *web.ArticleHandler {
	wire.Build(articleProviders)
	return &web.ArticleHandler{}
}

func InitArticleHandlerV1(articleDao dao.ArticleDao) *web.ArticleHandler {
	wire.Build(articleProvidersV1)
	return &web.ArticleHandler{}
}

func InitWebServer(templateId string) *gin.Engine {
	wire.Build(providers)
	return gin.Default()
}

func InitScheduler() *job.Scheduler {
	wire.Build(schedulerProvider)
	return new(job.Scheduler)
}
