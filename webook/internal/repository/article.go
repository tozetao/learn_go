package repository

import (
	"context"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
}

func NewArticleRepository(articleDao dao.ArticleDao) ArticleRepository {
	return &articleRepository{articleDao: articleDao}
}

type articleRepository struct {
	articleDao dao.ArticleDao
}

func (a *articleRepository) Create(ctx context.Context, article domain.Article) (int64, error) {
	return a.articleDao.Insert(ctx, dao.Article{
		Title:    article.Title,
		Content:  article.Content,
		AuthorID: article.Author.ID,
	})
}

func (a *articleRepository) Update(ctx context.Context, article domain.Article) error {
	// tip:
	// 用户只能更新自己的帖子。先查询再判定的性能不好，因为多了一次查询。
	// 正常用户是不会出现更新其他作者的帖子的，因为可以在更新时进行条件限制。

	return a.articleDao.UpdateByID(ctx, a.toEntity(article))
}

func (a *articleRepository) toEntity(article domain.Article) dao.Article {
	return dao.Article{
		Title:    article.Title,
		Content:  article.Content,
		ID:       article.ID,
		AuthorID: article.Author.ID,
	}
}
