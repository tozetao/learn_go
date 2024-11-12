package web

import (
	"github.com/gin-gonic/gin"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/service"
	"learn_go/webook/pkg/logger"
	"net/http"
)

type ArticleHandler struct {
	log logger.LoggerV2
	svc service.ArticleService
}

func NewArticleHandler(svc service.ArticleService, l logger.LoggerV2) *ArticleHandler {
	return &ArticleHandler{
		log: l,
		svc: svc,
	}
}

func (handler *ArticleHandler) Publish(c *gin.Context) {

	var req ArticleReq
	if c.Bind(&req) != nil {
		handler.log.Info("binding error during article publication.")
		return
	}

	claimsVal, _ := c.Get("user")
	userClaims, ok := claimsVal.(*UserClaims)
	if !ok {
		handler.log.Warn("get user claims error")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	articleID, err := handler.svc.Publish(c, req.toDomain(userClaims.Uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{
			Code: 5,
			Msg:  "failed",
		})
		return
	}

	c.JSON(http.StatusOK, Result{
		Msg:  "ok",
		Data: articleID,
	})
}

func (handler *ArticleHandler) Edit(c *gin.Context) {
	var req ArticleReq
	if c.Bind(&req) != nil {
		handler.log.Info("a binding error occurred while editing the article.")
		return
	}

	claimsVal, _ := c.Get("user")
	userClaims, ok := claimsVal.(*UserClaims)
	if !ok {
		handler.log.Warn("get user claims error")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	articleID, err := handler.svc.Save(c, req.toDomain(userClaims.Uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result{
			Code: 5,
			Msg:  "failed",
		})
		return
	}

	c.JSON(http.StatusOK, Result{
		Msg:  "ok",
		Data: articleID,
	})
}

func (handler *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")
	g.POST("/edit", handler.Edit)
	g.POST("/publish", handler.Publish)
}

type ArticleReq struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (req ArticleReq) toDomain(uid int64) domain.Article {
	return domain.Article{
		ID:      req.ID,
		Title:   req.Title,
		Content: req.Content,
		Author:  domain.Author{ID: uid},
	}
}
