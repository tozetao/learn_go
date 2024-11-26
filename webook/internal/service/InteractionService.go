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
}

func NewInteractionService(repo repository.InteractionRepository) InteractionService {
	return &interactionService{
		repo: repo,
	}
}

func (svc interactionService) Like(ctx context.Context, uid int64, articleID int64) error {
	//TODO implement me
	panic("implement me")
}

func (svc interactionService) View(ctx context.Context, articleID int64) error {
	return svc.repo.IncrReadCnt(ctx, articleID)
}
