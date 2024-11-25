package web

import (
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/service"
	"learn_go/webook/pkg/ginx"
	"learn_go/webook/pkg/logger"
	"net/http"
	"time"
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

func (handler *ArticleHandler) Publish(c *gin.Context, req ArticleReq, claims *UserClaims) (ginx.Result, error) {
	articleID, err := handler.svc.Publish(c, req.toDomain(claims.Uid))
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "failed",
		}, err
	}
	return ginx.Result{
		Msg:  "ok",
		Data: articleID,
	}, nil
}

func (handler *ArticleHandler) Withdraw(c *gin.Context) {
	type Req struct {
		ID int64 `json:"id"`
	}
	var req Req
	if c.Bind(&req) != nil {
		handler.log.Info("binding error during article publication.")
		return
	}

	claimsVal := c.MustGet("user")
	userClaims, ok := claimsVal.(*UserClaims)
	if !ok {
		handler.log.Warn("get user claims error")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err := handler.svc.Withdraw(c, domain.Article{
		ID: req.ID,
		Author: domain.Author{
			ID: userClaims.Uid,
		},
	})
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "failed",
		})
		return
	}

	c.JSON(http.StatusOK, Result{
		Msg: "ok",
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
		c.JSON(http.StatusOK, Result{
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

/*
创作者：

	查看自己的文章列表, /articles/list
	查看文章详情, /articles/detail

用户：

	从线上库查看已发布的文章
	/articles/published
*/

// List 创作者获取文章列表
func (handler *ArticleHandler) List(c *gin.Context, req ListReq, userClaims *UserClaims) (ginx.Result, error) {
	// offset, limit
	arts, err := handler.svc.GetList(c, userClaims.Uid, req.Offset, req.Limit)
	if err != nil {
		handler.log.Info("get list error during article publication.", logger.Error(err))
	}
	vos := slice.Map(arts, func(idx int, src domain.Article) ArticleVO {
		return ArticleVO{
			ID:      src.ID,
			Title:   src.Title,
			Content: src.Content,
			CTime:   src.CTime.Format(time.DateTime),
			UTime:   src.UTime.Format(time.DateTime),
		}
	})
	return ginx.Result{
		Msg:  "ok",
		Data: vos,
	}, nil
}

func (handler *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")
	g.POST("/edit", handler.Edit)
	g.POST("/publish", ginx.WrapBodyAndClaims(handler.Publish))

	// 查询作者的文章列表
	g.GET("/list", ginx.WrapBodyAndClaims(handler.List))
	// 查询作者的文章详情
	// g.GET("/detail", ginx.)
}
