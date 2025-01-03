package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

type App struct {
	server *gin.Engine
	//// 消费者服务
	//consumers []saramax.Consumer
	// cron
	cron *cron.Cron
}
