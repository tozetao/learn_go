package ioc

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	intrv1 "learn_go/webook/api/proto/gen/intr"
	"learn_go/webook/interaction/service"
	"learn_go/webook/internal/web/client"
)

// 构建grpc server、client

func NewGRPCInteractionServiceClient(service service.InteractionService) intrv1.InteractionServiceClient {
	type config struct {
		Addr      string
		Secure    bool
		Threshold int32
	}

	var cfg config
	err := viper.UnmarshalKey("grpc.client.addr", &cfg)
	if err != nil {
		panic(err)
	}
	options := make([]grpc.DialOption, 4)
	if cfg.Secure {
		// 添加https相关配置
	} else {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	cc, err := grpc.Dial(cfg.Addr, options...)
	if err != nil {
		panic(err)
	}

	// rpc
	remote := intrv1.NewInteractionServiceClient(cc)

	local := client.NewInteractionServiceAdapter(service)
	interSvcClient := client.NewInteractionServiceClient(local, remote, cfg.Threshold)

	// 当配置文件变动时重新加载配置
	viper.OnConfigChange(func(e fsnotify.Event) {
		interSvcClient.UpdateThreshold(cfg.Threshold)
	})
	return interSvcClient
}
