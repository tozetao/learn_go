//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"learn_go/webook/interaction/event/article"
	grpc2 "learn_go/webook/interaction/grpc"
	"learn_go/webook/interaction/ioc"
	"learn_go/webook/interaction/repository"
	"learn_go/webook/interaction/repository/cache"
	"learn_go/webook/interaction/repository/dao"
	"learn_go/webook/interaction/service"
)

// 第三方依赖
var thirdPartySet = wire.NewSet(
	ioc.NewLogger,
	ioc.NewDB,
	ioc.NewRedis,

	ioc.NewSaramaConfig,
	ioc.NewConsumerClient,
)

var interactionSvcSet = wire.NewSet(
	service.NewInteractionService,
	repository.NewInteractionRepository,
	dao.NewInteractionDao,
	cache.NewInteractionCache,
)

func InitApp() *App {
	wire.Build(thirdPartySet,
		interactionSvcSet,
		grpc2.NewInteractionServiceServer,

		article.NewBatchReadEventConsumer,
		ioc.NewConsumers,

		ioc.InitGRPCServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
