package main

import (
	"learn_go/webook/pkg/grpcx"
	"learn_go/webook/pkg/saramax"
)

type App struct {
	// 消费者服务
	consumers []saramax.Consumer

	server *grpcx.Server
}
