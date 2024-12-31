//go:build wireinject
// +build wireinject

package startup

import (
	"github.com/google/wire"
	"learn_go/webook/interaction/grpc"
	"learn_go/webook/interaction/repository"
	"learn_go/webook/interaction/repository/cache"
	"learn_go/webook/interaction/repository/dao"
	"learn_go/webook/interaction/service"
	"learn_go/webook/ioc"
)

var thirdPartySet = wire.NewSet(
	NewDB,
	NewRedis,
	ioc.NewLogger,
)

var interactionSet = wire.NewSet(
	service.NewInteractionService,
	repository.NewInteractionRepository,
	dao.NewInteractionDao,
	cache.NewInteractionCache,
)

func InitInteractionService() service.InteractionService {
	wire.Build(thirdPartySet, interactionSet)
	return service.NewInteractionService(nil)
}

func InitInteractionServiceServer() *grpc.InteractionServiceServer {
	wire.Build(thirdPartySet, interactionSet, grpc.NewInteractionServiceServer)
	return &grpc.InteractionServiceServer{}
}
