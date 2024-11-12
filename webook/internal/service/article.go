package service

import (
	"context"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/repository/article"
	"learn_go/webook/pkg/logger"
)

/*
Publish的数据同步。

发布这种操作会有3种情况：
1. 新建Article并发表：插入制作库 => 插入线上库
2. 编辑Article后首次发表：更新制作库 => 插入线上库
3. 编辑Article并发表：更新制作库 => 更新线上库

可以发现制作库的数据需要同步到线上库，在代码层面上可以着由划分：
ArticleAuthorService
	Save(domain.Article)
ArticleReaderService
	Save(domain.Article)
偏向微服务的划分。根据线上库、制作库划分为俩个不同的服务，最后在web层这俩种服务。


在Repository层分为ArticleAuthorRepository、ArticleReaderRepository，
在Service层聚合这俩个仓库，实现发布。


将数据同步的操作封装到ArticleRepository种。
ArticleRepository
	Save(domain.Article)


*/

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)

	Publish(ctx context.Context, article domain.Article) (int64, error)

	PublishV1(ctx context.Context, article domain.Article) (int64, error)
}

type articleService struct {
	articleRepo article.ArticleRepository
	log         logger.LoggerV2
}

func (svc *articleService) Publish(ctx context.Context, article domain.Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (svc *articleService) PublishV1(ctx context.Context, article domain.Article) (int64, error) {

}

func NewArticleService(articleRepo article.ArticleRepository, log logger.LoggerV2) ArticleService {
	return &articleService{
		log:         log,
		articleRepo: articleRepo,
	}
}

func (svc *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	if article.ID > 0 {
		err := svc.articleRepo.Update(ctx, article)
		return article.ID, err
	}
	return svc.articleRepo.Create(ctx, article)
}
