package article

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/repository"
	"learn_go/webook/internal/repository/cache"
	"learn_go/webook/internal/repository/dao"
	"learn_go/webook/pkg/logger"
	"time"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error

	Sync(ctx context.Context, article domain.Article) (int64, error)

	SyncStatus(ctx context.Context, id int64, authorID int64, status int8) error

	// SyncV1 在repository层同步数据
	SyncV1(ctx context.Context, article domain.Article) (int64, error)

	GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetByID(ctx context.Context, id int64) (domain.Article, error)
	ListPub(ctx context.Context, t time.Time, offset int, limit int) ([]domain.Article, error)
	GetPubByID(ctx context.Context, id int64) (domain.Article, error)
}

func NewArticleRepository(articleDao dao.ArticleDao, articleCache cache.ArticleCache, userRepo repository.UserRepository, log logger.LoggerV2) ArticleRepository {
	return &articleRepository{
		articleDao:   articleDao,
		log:          log,
		articleCache: articleCache,
		userRepo:     userRepo,
	}
}

type articleRepository struct {
	articleDao dao.ArticleDao

	articleAuthorDao dao.ArticleAuthorDao
	articleReaderDao dao.ArticleReaderDao

	articleCache cache.ArticleCache

	userRepo repository.UserRepository

	log logger.LoggerV2
}

func (repo *articleRepository) GetPubByID(ctx context.Context, id int64) (domain.Article, error) {
	// 1. 先从缓存中读取
	res, err := repo.articleCache.GetPub(ctx, id)
	if err == nil {
		return res, nil
	}
	// 2. 再从数据库中读取
	pubArt, err := repo.articleDao.GetPubByID(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	// 2.1 查询用户信息
	user, err := repo.userRepo.FindByID(ctx, pubArt.AuthorID)
	if err != nil {
		return domain.Article{}, err
	}
	article := repo.toDomain(dao.Article(pubArt))
	article.Author.Name = user.Nickname

	// 3. 重新载入缓存
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		err := repo.articleCache.SetPub(ctx, article)
		if err != nil {
			// 记录日志
		}
	}()

	return article, nil
}

func (repo *articleRepository) ListPub(ctx context.Context, t time.Time, offset int, limit int) ([]domain.Article, error) {
	arts, err := repo.articleDao.ListPub(ctx, t, offset, limit)
	if err != nil {
		return nil, err
	}
	newArts := slice.Map[dao.Article, domain.Article](arts, func(idx int, src dao.Article) domain.Article {
		return repo.toDomain(src)
	})
	return newArts, nil
}

// GetByAuthor
// 如何实现缓存?
// 因为offset、limit是可变的，这个接口很难做缓存，因此我们已uid作为key，只缓存作者的第一页数据。
// 什么时候清除缓存?
// Create、Update、Sync，当作者执行这三个方法时，需要清除缓存。
func (repo *articleRepository) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	if offset == 0 && limit <= 100 {
		res, err := repo.articleCache.GetFirstPage(ctx, uid)
		if err == nil {
			return res, nil
		} else {
			// 要考虑记录日志
			// 缓存未命中是可以忽略的

			// 是否需要记录非redis.Nil的错误?
			//if err != cache.ErrKeyNotExists {
			//	repo.log.Error("获取第一张文章错误", logger.Error(err))
			//}
		}
	}
	arts, err := repo.articleDao.GetByAuthor(ctx, uid, offset, limit)
	if err != nil {
		return nil, err
	}
	res := slice.Map(arts, func(idx int, src dao.Article) domain.Article {
		return repo.toDomain(src)
	})

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
		defer cancel()

		if offset == 0 && limit <= 100 {
			err = repo.articleCache.SetFirstPage(ctx, uid, res)
			// 记录错误日志并上报
			if err != nil {
				repo.log.Error("缓存文章列表失败", logger.Error(err))
			}
		}
	}()

	// 预加载
	// 指的是提现加载用户可能在近期关联操作的数据。比如作者查看文章列表后，可能会立即点击第一篇文章，这时候可以提前缓存第一篇文章。
	// 注：这是一种用户行为上测预测，所以缓存命中率可能不高，因此过期时间要尽可能短
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		if offset == 0 && limit <= 100 {
			repo.preCache(ctx, res)
		}
	}()

	return res, nil
}

func (repo *articleRepository) GetByID(ctx context.Context, id int64) (domain.Article, error) {
	// 1. 先从缓存中查询
	article, err := repo.articleCache.Get(ctx, id)
	if err == nil {
		return article, nil
	}

	// 2. 没有则从数据库中查询
	entity, err := repo.articleDao.GetByID(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	article = repo.toDomain(entity)

	// 3. 再载入缓存
	go func() {
		err := repo.articleCache.Set(ctx, article)
		if err != nil {
			repo.log.Error("缓存作者文章失败", logger.Error(err))
		}
	}()

	return article, nil
}

func (repo *articleRepository) preCache(ctx context.Context, arts []domain.Article) {
	if len(arts) > 0 {
		err := repo.articleCache.Set(ctx, arts[0])
		if err != nil {
			repo.log.Error("预缓存失败", logger.Error(err))
		}
	}
}

func (repo *articleRepository) Sync(ctx context.Context, article domain.Article) (int64, error) {
	// 如果是同个库，可以把数据的同步交给dao层。
	id, err := repo.articleDao.Sync(ctx, repo.toEntity(article))
	if err != nil {
		return 0, err
	}

	// 删除缓存
	err = repo.articleCache.RemoveFirstPage(ctx, id)
	if err != nil {
		repo.log.Error("从缓存中移除文章列表失败", logger.Error(err))
	}

	go func() {
		// 尝试重新载入发布的文章
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		// 查询用户信息
		user, err := repo.userRepo.FindByID(ctx, article.Author.ID)
		if err != nil {
			// 记录日志
		}
		article.Author = domain.Author{
			ID:   user.ID,
			Name: user.Nickname,
		}
		err = repo.articleCache.SetPub(ctx, article)
		if err != nil {
			// 记录日志
		}
	}()

	return id, err
}

func (repo *articleRepository) SyncV1(ctx context.Context, article domain.Article) (int64, error) {
	// 数据可能是俩个存储源（比如不同库），就由不同存储源所对应的dao来处理。
	var (
		id  = article.ID
		err error
	)
	if article.ID > 0 {
		err = repo.articleAuthorDao.UpdateByID(ctx, repo.toEntity(article))
	} else {
		id, err = repo.articleAuthorDao.Insert(ctx, repo.toEntity(article))
	}
	if err != nil {
		return 0, err
	}
	article.ID = id
	err = repo.articleReaderDao.Upsert(ctx, repo.toEntity(article))
	return id, err
}

func (repo *articleRepository) SyncStatus(ctx context.Context, id int64, authorID int64, status int8) error {
	return repo.articleDao.SyncStatus(ctx, id, authorID, status)
}

func (repo *articleRepository) Create(ctx context.Context, article domain.Article) (int64, error) {
	id, err := repo.articleDao.Insert(ctx, repo.toEntity(article))
	if err == nil {
		// 删除缓存
		go func() {
			err := repo.articleCache.RemoveFirstPage(ctx, article.Author.ID)
			if err != nil {
				repo.log.Error("从缓存中移除文章列表失败", logger.Error(err))
			}
		}()
	}
	return id, err
}

func (repo *articleRepository) Update(ctx context.Context, article domain.Article) error {
	// tip:
	// 用户只能更新自己的帖子。先查询再判定的性能不好，因为多了一次查询。
	// 正常用户是不会出现更新其他作者的帖子的，因为可以在更新时进行条件限制。
	err := repo.articleDao.UpdateByID(ctx, repo.toEntity(article))

	if err == nil {
		// 删除缓存
		go func() {
			err := repo.articleCache.RemoveFirstPage(ctx, article.Author.ID)
			if err != nil {
				repo.log.Error("从缓存中移除文章列表失败", logger.Error(err))
			}
		}()

	}
	return err
}

func (repo *articleRepository) toDomain(src dao.Article) domain.Article {
	return domain.Article{
		ID:      src.ID,
		Title:   src.Title,
		Content: src.Content,
		Author: domain.Author{
			ID: src.AuthorID,
		},
		Status: domain.ArticleStatus(src.Status),
		CTime:  time.UnixMilli(src.Ctime),
		UTime:  time.UnixMilli(src.Utime),
	}
}

func (repo *articleRepository) toEntity(article domain.Article) dao.Article {
	return dao.Article{
		Title:    article.Title,
		Content:  article.Content,
		ID:       article.ID,
		AuthorID: article.Author.ID,
		Status:   article.Status.ToInt8(),
	}
}
