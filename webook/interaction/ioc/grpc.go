package ioc

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	intrv1 "learn_go/webook/api/proto/gen/intr"
	grpc2 "learn_go/webook/interaction/grpc"
	"learn_go/webook/pkg/grpcx"
)

func InitGRPCServer(intrSvcServer *grpc2.InteractionServiceServer) *grpcx.Server {
	type config struct {
		Addr string
	}
	var cfg config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	intrv1.RegisterInteractionServiceServer(server, intrSvcServer)

	return &grpcx.Server{
		Server: server,
		Addr:   cfg.Addr,
	}
}
