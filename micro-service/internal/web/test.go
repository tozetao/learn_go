package web

import (
	"github.com/gin-gonic/gin"
	"learn_go/webook/pkg/logger"
	"net/http"
)

type TestHandler struct {
	logger logger.LoggerV2
}

func NewTestHandler(logger logger.LoggerV2) *TestHandler {
	return &TestHandler{
		logger: logger,
	}
}

func (h *TestHandler) Test(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

func (h *TestHandler) RegisterHandler(server *gin.Engine) {
	server.GET("/test", h.Test)
}
