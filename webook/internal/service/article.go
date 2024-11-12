package service

import (
	"context"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/repository"
	"learn_go/webook/pkg/logger"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)

	Publish(ctx context.Context, article domain.Article) (int64, error)
}

type articleService struct {
	articleRepo repository.ArticleRepository
	log         logger.LoggerV2
}

func (a articleService) Publish(ctx context.Context, article domain.Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func NewArticleService(articleRepo repository.ArticleRepository, log logger.LoggerV2) ArticleService {
	return &articleService{
		log:         log,
		articleRepo: articleRepo,
	}
}

func (a articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	if article.ID > 0 {
		err := a.articleRepo.Update(ctx, article)
		return article.ID, err
	}
	return a.articleRepo.Create(ctx, article)
}
