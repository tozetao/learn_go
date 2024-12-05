package web

import (
	"errors"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/service"
	"learn_go/webook/pkg/ginx"
	"learn_go/webook/pkg/logger"
	"net/http"
	"strconv"
	"time"
)

type ArticleHandler struct {
	log      logger.LoggerV2
	svc      service.ArticleService
	interSvc service.InteractionService

	biz string
}

func NewArticleHandler(svc service.ArticleService, interSvc service.InteractionService, l logger.LoggerV2) *ArticleHandler {
	return &ArticleHandler{
		log:      l,
		svc:      svc,
		interSvc: interSvc,
		biz:      "article",
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

// List 创作者获取文章列表
func (handler *ArticleHandler) List(c *gin.Context, req ListReq, userClaims *UserClaims) (ginx.Result, error) {
	// offset, limit
	arts, err := handler.svc.GetList(c, userClaims.Uid, req.Offset, req.Limit)
	if err != nil {
		handler.log.Info("get list error during article publication.", logger.Error(err))
	}
	vos := slice.Map(arts, func(idx int, src domain.Article) ArticleVO {
		return handler.ToVO(src)
	})
	return ginx.Result{
		Msg:  "ok",
		Data: vos,
	}, nil
}

func (handler *ArticleHandler) Detail(c *gin.Context) {
	// 文章ID、作者ID
	idstr := c.Param("id")
	articleID, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		c.JSON(200, ginx.Result{
			Code: 5,
			Msg:  "params error",
		})
		handler.log.Warn("解析id出错", logger.Error(err))
		return
	}

	art, err := handler.svc.GetByID(c, articleID)
	if err != nil {
		c.JSON(200, ginx.Result{
			Code: 5,
			Msg:  "failed",
		})
		//记录日志
		return
	}
	claims, ok := c.MustGet("user").(*UserClaims)
	if !ok {
		c.JSON(200, ginx.Result{
			Code: 5,
			Msg:  "auth error",
		})
		handler.log.Warn("get user claims error")
		return
	}
	if claims.Uid != art.Author.ID {
		c.JSON(200, ginx.Result{
			Code: 5,
			Msg:  "auth error",
		})
		// 记录日志
		return
	}
	c.JSON(200, ginx.Result{
		Msg:  "ok",
		Data: handler.ToVO(art),
	})
}

// GetPublished 获取已发布的文章
func (handler *ArticleHandler) GetPublished(c *gin.Context) {
	startStr := c.Query("start")
	loc, err := time.LoadLocation("Asis/Shanghai")
	if err != nil {
		c.JSON(200, ginx.Result{Code: 5, Msg: "parse time location error"})
		return
	}
	start, err := time.ParseInLocation(time.DateTime, startStr, loc)
	if err != nil {
		c.JSON(200, ginx.Result{Code: 5, Msg: "parse time error"})
		return
	}

	offsetStr := c.Query("offset")
	offset, err := strconv.ParseInt(offsetStr, 10, 0)
	if err != nil {
		c.JSON(200, ginx.Result{Code: 5, Msg: "params error"})
		return
	}
	limitStr := c.Query("limit")
	limit, err := strconv.ParseInt(limitStr, 10, 0)
	if err != nil {
		c.JSON(200, ginx.Result{Code: 5, Msg: "params error"})
		return
	}

	arts, err := handler.svc.ListPub(c, start, int(offset), int(limit))
	if err != nil {
		c.JSON(200, ginx.Result{
			Code: 5,
			Msg:  "internal server error",
		})
		return
	}
	c.JSON(200, ginx.Result{
		Msg: "ok",
		Data: slice.Map[domain.Article, ArticleVO](arts, func(idx int, src domain.Article) ArticleVO {
			return handler.ToVO(src)
		}),
	})
}

func (handler *ArticleHandler) PubDetail(c *gin.Context) {
	idStr := c.Param("id")
	articleID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(200, ginx.Result{Code: 5, Msg: "params error"})
		return
	}

	userClaims, ok := c.MustGet("user").(*UserClaims)
	if !ok {
		c.JSON(200, ginx.Result{Code: 5, Msg: "auth error"})
		return
	}

	art, err := handler.svc.GetPubArticle(c, articleID)
	if err != nil {
		c.JSON(200, ginx.Result{Code: 5, Msg: "internal server error."})
		return
	}
	vo := handler.ToVO(art)

	// 查询文章的点赞数、收藏数和观看书
	inter, err := handler.interSvc.Get(c, userClaims.Uid, handler.biz, articleID)
	if err != nil {
		// 记录错误
	} else {
		// 查询用户是否点赞、收藏
		vo.Views = inter.Views
		vo.Favorites = inter.Favorites
		vo.Likes = inter.Likes
		vo.Liked = inter.Liked
		vo.Collected = inter.Collected
	}
	//ok, _ = handler.interSvc.Liked(c, userClaims.Uid, handler.biz, articleID)
	//if !ok {
	//	vo.Liked = ArticleLike
	//} else {
	//	vo.Liked = ArticleUnlike
	//}

	// 增加阅读数
	go func() {
		err := handler.interSvc.View(c, handler.biz, art.ID)
		if err != nil {
			// 只能记录日志，上传告警信息
		}
	}()

	c.JSON(200, ginx.Result{
		Msg:  "ok",
		Data: vo,
	})
}

func (handler *ArticleHandler) Like(c *gin.Context, req LikeReq, userClaims *UserClaims) (ginx.Result, error) {
	var err error
	if req.Action == ArticleLike {
		err = handler.interSvc.Like(c, userClaims.Uid, handler.biz, req.ArticleID)
	} else {
		err = handler.interSvc.CancelLike(c, userClaims.Uid, handler.biz, req.ArticleID)
	}
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "failed",
		}, errors.New("点赞或取消点赞失败")
	}
	return ginx.Result{
		Msg: "success",
	}, nil
}

func (handler *ArticleHandler) Favorite(ctx *gin.Context, req FavoriteReq, userClaims *UserClaims) (ginx.Result, error) {
	err := handler.interSvc.Favorite(ctx, userClaims.Uid, req.FavoriteID, handler.biz, req.ArticleID)
	if err != nil {
		return ginx.Result{Code: 5, Msg: "save error"}, errors.New("failed")
	}
	return ginx.Result{Msg: "ok"}, nil
}

func (handler *ArticleHandler) ToVO(src domain.Article) ArticleVO {
	return ArticleVO{
		ID:      src.ID,
		Title:   src.Title,
		Content: src.Content,
		CTime:   src.CTime.Format(time.DateTime),
		UTime:   src.UTime.Format(time.DateTime),
	}
}

func (handler *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")
	g.POST("/edit", handler.Edit)
	g.POST("/publish", ginx.WrapBodyAndClaims(handler.Publish))

	// 查询作者的文章列表
	g.GET("/list", ginx.WrapBodyAndClaims(handler.List))
	// 查询作者的文章详情
	g.GET("/detail/:id", handler.Detail)

	// 已发布文章接口
	pub := g.Group("/pub")
	pub.GET("/details/:id", handler.PubDetail)

	// 点赞接口
	g.POST("/like", ginx.WrapBodyAndClaims[LikeReq, *UserClaims](handler.Like))

	g.POST("/favorite", ginx.WrapBodyAndClaims[FavoriteReq, *UserClaims](handler.Favorite))
}

/*

接下来实现文章的 阅读量、点赞、收藏

article_interaction（文章互动）
	id,UTime,CTime,read_cnt, likes,

article_like

article_favorite

文章阅读
	当用户查看一篇文章时，增加文章的阅读量。




点赞
	用户点赞一篇文章，文章点赞量+1，用户点赞列表+1
取消点赞
	文章点赞量-1，用户点赞列表-1

收藏
	文章的收藏数+1，用户的收藏夹+1
取消收藏

架构、代码结构
如果我们的应用是一个大型应用，采用了微服务架构，那么阅读量、点赞、收藏确实是可以分成3个单独的服务。
但是对于单体应用，小应用，我们可以将这3个聚合在一起。

微服务：
// 文章读数计数服务
ArticleReadCntService

// 文章点赞服务
ArticleLikeService

// 文章收藏服务
ArticleFavorite

单体服务：
InteractionService
	view
	like
	favorite

InteractionRepository

InteractionDao

Interaction
	id, biz, biz_id, c_time, u_time, likes, favorites, views
	读取性能高，写入性能差。
	当读取一篇文章时可以一次性的将三个指标都读取出来，但是触发点赞、收藏等行为时，由于要增加计数，在update时会加锁等待。

	id, biz, biz_id, type,
	通过type来区分不同的指标，因此写入性能会好些，但是读取差，因为同个资源的多个指标需要从磁盘上随机读，无法顺序读（比如查看后再点赞，俩条记录会相差的比较远）




收藏接口的实现

文章的交互数据的获取
	交互数据的获取可以有俩种方式：
	1. 将Interaction领域对象作为Article领域对象的值对象，从概念上来说，这些交互数据的获取要在Article Repository层面上来解决问题。
	2. 新建一个service，聚合article、interaction服务来完成数据的组装。我们这里直接在handler层面上来聚合服务了。


interactionService
	Liked
	collected

	Get(): 获取文章的点赞数、收藏数、阅读数
	为Interaction领域添加是否点赞、是否收藏属性。那么从语义上来说，应该Repository层面上获取这俩个属性的值。


*/
