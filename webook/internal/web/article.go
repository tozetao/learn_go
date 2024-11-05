package web

import "github.com/gin-gonic/gin"

type ArticleHandler struct {
}

func (h *ArticleHandler) Edit(c *gin.Context) {

}

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")
	g.POST("/edit", h.Edit)
}
