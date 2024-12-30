package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"learn_go/webook/internal/domain"
	"strings"
	"time"
)

type ArticleCache interface {
	GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, uid int64, arts []domain.Article) error
	RemoveFirstPage(ctx context.Context, id int64) error

	Set(ctx context.Context, article domain.Article) error
	Get(ctx context.Context, id int64) (domain.Article, error)
	GetPub(ctx context.Context, id int64) (domain.Article, error)
	SetPub(ctx context.Context, article domain.Article) error
}

type articleCache struct {
	cmd redis.Cmdable
}

func NewArticleCache(cmd redis.Cmdable) ArticleCache {
	return &articleCache{
		cmd: cmd,
	}
}

func (cache *articleCache) listKey(uid int64) string {
	return fmt.Sprintf("article:list:%d", uid)
}

func (cache *articleCache) articleKey(articleID int64) string {
	return fmt.Sprintf("article:detail:%d", articleID)
}

func (cache *articleCache) pubArticleKey(articleID int64) string {
	return fmt.Sprintf("pub_article:detail:%d", articleID)
}

func (cache *articleCache) SetPub(ctx context.Context, article domain.Article) error {
	data, err := json.Marshal(article)
	if err != nil {
		return err
	}
	return cache.cmd.Set(ctx, cache.pubArticleKey(article.ID), data, time.Minute*2).Err()
}

func (cache *articleCache) GetPub(ctx context.Context, id int64) (domain.Article, error) {
	data, err := cache.cmd.Get(ctx, cache.pubArticleKey(id)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var article domain.Article
	err = json.Unmarshal(data, &article)
	if err != nil {
		return domain.Article{}, err
	}
	return article, nil
}

func (cache *articleCache) Set(ctx context.Context, article domain.Article) error {
	data, err := json.Marshal(article)
	if err != nil {
		return err
	}
	return cache.cmd.Set(ctx, cache.articleKey(article.ID), data, time.Minute*2).Err()
}

func (cache *articleCache) Get(ctx context.Context, id int64) (domain.Article, error) {
	data, err := cache.cmd.Get(ctx, cache.articleKey(id)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var article domain.Article
	err = json.Unmarshal(data, &article)
	if err != nil {
		return domain.Article{}, err
	}
	return article, nil
}

func (cache *articleCache) GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error) {
	data, err := cache.cmd.Get(ctx, cache.listKey(uid)).Result()
	if err != nil {
		return nil, err
	}

	var arts []domain.Article
	err = json.NewDecoder(strings.NewReader(data)).Decode(&arts)
	if err != nil {
		return nil, err
	}
	return arts, nil
}

func (cache *articleCache) SetFirstPage(ctx context.Context, uid int64, arts []domain.Article) error {
	artsJson, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	return cache.cmd.Set(ctx, cache.listKey(uid), artsJson, time.Minute*5).Err()
}

func (cache *articleCache) RemoveFirstPage(ctx context.Context, id int64) error {
	return cache.cmd.Del(ctx, cache.listKey(id)).Err()
}
