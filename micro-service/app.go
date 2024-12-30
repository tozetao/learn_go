package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"learn_go/webook/internal/event"
)

type App struct {
	server *gin.Engine
	// 消费者服务
	consumers []event.Consumer
	// cron
	cron *cron.Cron
}
