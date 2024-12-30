package article

import (
	"context"
	"learn_go/webook/internal/domain"
)

type AuthorRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
}

func (repo *articleAuthorRepository) Create(ctx context.Context, article domain.Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *articleAuthorRepository) Update(ctx context.Context, article domain.Article) error {
	//TODO implement me
	panic("implement me")
}

type articleAuthorRepository struct {
}

func NewArticleAuthorRepository() AuthorRepository {
	return &articleAuthorRepository{}
}
