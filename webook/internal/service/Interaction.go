package service

import (
	"context"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/repository"
)

type InteractionService interface {
	View(ctx context.Context, articleID int64) error

	Like(ctx context.Context, uid int64, articleID int64) error
	CancelLike(ctx context.Context, uid int64, articleID int64) error
	Favorite(ctx context.Context, uid int64, favoriteID int64, articleID int64) error

	Get(ctx context.Context, biz string, bizID int64) (domain.Interaction, error)

	Liked(ctx context.Context, uid int64, biz string, bizID int64) bool
	Collected(ctx context.Context, uid int64, biz string, bizID int64) bool
}

func (svc *interactionService) Liked(ctx context.Context, uid int64, biz string, bizID int64) bool {
	userLike, err := svc.repo.GetUserLikeInfo(ctx, uid, biz, bizID)
	if err != nil {
		// 记录日志，报警。当然错误可能是ErrRecordNotFound。我们整个项目中，repository找不到记录都是以错误返回的。
		return false
	}
	return userLike.Uid == uid && userLike.BizID == bizID
}

func (svc *interactionService) Collected(ctx context.Context, uid int64, biz string, bizID int64) bool {
	// TODO: implement me
	return false
}

func (svc *interactionService) Favorite(ctx context.Context, uid int64, favoriteID int64, articleID int64) error {
	return svc.repo.AddFavorite(ctx, uid, favoriteID, svc.biz, articleID)
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

func (svc *interactionService) Get(ctx context.Context, biz string, bizID int64) (domain.Interaction, error) {
	return svc.repo.Get(ctx, biz, bizID)
}

func (svc *interactionService) Like(ctx context.Context, uid int64, articleID int64) error {
	return svc.repo.IncrLike(ctx, uid, svc.biz, articleID)
}

func (svc *interactionService) CancelLike(ctx context.Context, uid int64, articleID int64) error {
	return svc.repo.DecrLike(ctx, uid, svc.biz, articleID)
}

// View 查看文章：增加文章点击量
func (svc *interactionService) View(ctx context.Context, articleID int64) error {
	return svc.repo.IncrReadCnt(ctx, svc.biz, articleID)
}
