package main

import (
	"github.com/gin-gonic/gin"
	"learn_go/webook/internal/event"
)

type App struct {
	server *gin.Engine
	// 消费者服务
	consumers []event.Consumer
}
