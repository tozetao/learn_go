//go:build wireinject
// +build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"learn_go/webook/internal/repository"
	"learn_go/webook/internal/repository/article"
	"learn_go/webook/internal/repository/cache"
	"learn_go/webook/internal/repository/dao"
	"learn_go/webook/internal/service"
	"learn_go/webook/internal/web"
	"learn_go/webook/ioc"
)

var (
	providers = wire.NewSet(
		// 第三方依赖
		ioc.NewLogger,
		NewDB,
		NewRedis,

		cache.NewCodeCache, cache.NewUserCache,

		dao.NewUserDao,

		repository.NewCodeRepository,
		repository.NewUserRepository,

		ioc.InitSMSService,
		ioc.InitOAuth2Service,
		service.NewCodeService,
		service.NewUserService,

		web.NewSMSHandler,
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		web.NewJWTHandler,
		web.NewTestHandler,

		ioc.InitMiddlewares,
		ioc.InitGin,
	)

	articleProviders = wire.NewSet(
		ioc.NewLogger,
		NewDB,
		NewRedis,
		web.NewArticleHandler,

		service.NewArticleService,
		article.NewArticleRepository,
		article.NewArticleAuthorRepository,
		article.NewArticleReaderRepository,
		dao.NewArticleDao,
	)

	articleProvidersV1 = wire.NewSet(
		NewDB,
		NewRedis,
		NewMangoDB,
		ioc.NewLogger,
		web.NewArticleHandler,

		service.NewArticleService,
		article.NewArticleRepository,
		article.NewArticleAuthorRepository,
		article.NewArticleReaderRepository,
	)
)

func InitArticleHandler() *web.ArticleHandler {
	wire.Build(articleProviders)
	return &web.ArticleHandler{}
}

func InitArticleHandlerV1(articleDao dao.ArticleDao) *web.ArticleHandler {
	wire.Build(articleProvidersV1)
	return &web.ArticleHandler{}
}

func InitWebServer(templateId string) *gin.Engine {
	wire.Build(providers)
	return gin.Default()
}
