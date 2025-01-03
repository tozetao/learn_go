//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	repository2 "learn_go/webook/interaction/repository"
	cache2 "learn_go/webook/interaction/repository/cache"
	dao2 "learn_go/webook/interaction/repository/dao"
	service2 "learn_go/webook/interaction/service"
	event "learn_go/webook/internal/event/article"
	"learn_go/webook/internal/repository"
	"learn_go/webook/internal/repository/article"
	"learn_go/webook/internal/repository/cache"
	"learn_go/webook/internal/repository/dao"
	"learn_go/webook/internal/service"
	"learn_go/webook/internal/web"
	"learn_go/webook/ioc"
)

var rankingSet = wire.NewSet(
	service.NewRankingService,
	repository.NewRankingRepository,
	ioc.NewRedisRanking,
	ioc.NewLocalCacheRanking,
)

// 第三方依赖
var thirdPartySet = wire.NewSet(
	ioc.NewLogger,
	ioc.NewDB,
	ioc.NewRedis,
	ioc.InitMiddlewares,
	ioc.InitGin,
)

var jobSet = wire.NewSet(
	ioc.InitRankingJob,
	ioc.InitCron,
)

// 生产者
var producerSet = wire.NewSet(
	ioc.NewSaramaConfig,
	ioc.NewSyncProducer,
	event.NewSyncProducer,
)

var articleSet = wire.NewSet(
	web.NewArticleHandler,
	service.NewArticleService,

	article.NewArticleRepository,
	article.NewArticleAuthorRepository,
	article.NewArticleReaderRepository,
	dao.NewArticleDao,
	cache.NewArticleCache,

	service2.NewInteractionService,
	repository2.NewInteractionRepository,
	dao2.NewInteractionDao,
	cache2.NewInteractionCache,

	ioc.NewGRPCInteractionServiceClient,
)

var smsSet = wire.NewSet(
	web.NewSMSHandler,

	service.NewCodeService,
	ioc.NewSMSService,

	repository.NewCodeRepository,
	cache.NewCodeCache,
)

var userSet = wire.NewSet(
	web.NewUserHandler,
	service.NewUserService,
	repository.NewUserRepository,
	cache.NewUserCache,
	dao.NewUserDao,
)

var wechatSet = wire.NewSet(
	web.NewOAuth2WechatHandler,
	ioc.InitOAuth2Service,
	web.NewJWTHandler,
)

var (
	providers = wire.NewSet(
		thirdPartySet,
		producerSet,
		jobSet,

		rankingSet,
		articleSet,
		smsSet,
		userSet,
		wechatSet,

		web.NewTestHandler,
		wire.Struct(new(App), "*"),
	)
)

func InitApp(templateId string) *App {
	wire.Build(providers)
	return new(App)
}
