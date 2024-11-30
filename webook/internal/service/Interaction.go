package service

import (
	"context"
	"learn_go/webook/internal/repository"
)

type InteractionService interface {
	View(ctx context.Context, articleID int64) error

	Like(ctx context.Context, uid int64, articleID int64) error
}
type interactionService struct {
	repo repository.InteractionRepository
	biz  string
}

func NewInteractionService(repo repository.InteractionRepository) InteractionService {
	return &interactionService{
		repo: repo,
		biz:  "article",
	}
}

func (svc interactionService) Like(ctx context.Context, uid int64, articleID int64) error {

}

// View 查看文章：增加文章点击量
func (svc interactionService) View(ctx context.Context, articleID int64) error {
	return svc.repo.IncrReadCnt(ctx, svc.biz, articleID)
}
