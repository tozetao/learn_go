package article

import (
	"context"
	"learn_go/webook/internal/domain"
)

type ReaderRepository interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
}

type articleReaderRepository struct {
}

func (repo *articleReaderRepository) Save(ctx context.Context, article domain.Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func NewArticleReaderRepository() ReaderRepository {
	return &articleReaderRepository{}
}
