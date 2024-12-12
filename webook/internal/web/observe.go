package web

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

type ObserveHandler struct {
}

func (h *ObserveHandler) RegisterHandler(server *gin.Engine) {
	g := server.Group("/observe")
	g.GET("/metric", func(c *gin.Context) {
		sleep := rand.Int31n(1000)
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		c.String(http.StatusOK, "ok")
	})
}
