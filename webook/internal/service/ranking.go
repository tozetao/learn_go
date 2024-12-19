package service

import (
	"context"
	"github.com/ecodeclub/ekit/queue"
	"github.com/ecodeclub/ekit/slice"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/repository"
	"math"
	"time"
)

/*
如何保证高性能?
	要保证高性能必然会引入本地缓存。基本解决方案都是：本地缓存 + redis缓存 + mysql方案

操作
	查找：先找本地缓存，再查找redis，最后查找mysql
	更新：先更新数据库，再更新本地缓存，最后更新redis。

	为什么更新要先更新本地缓存，再更新redis?
		因为本地缓存基本不可能失败，当进程已经再处理请求了，更新本地缓存只是CPU操作，没有其他io操作。


如何保证高可用?
	整个榜单的功能都是依赖于redis和mysql。
	如果mysql奔溃了，定时任务肯定更新失败。
	如果redis奔溃了，当本地缓存失效后，redis也不可用，那么查询接口也肯定失败。

	规避mysql不可用的问题：
		让redis的缓存永不过期，如果mysql不可用了，即使定时任务无法执行也仍然会有数据可访问。


	从代码层面上做一个兜底：
		给Get接口做一个兜底，正常情况下是从本地缓存中获取，获取不到去redis里面获取。
		那么我们可以在redis奔溃的时候，再次从本地缓存获取，此时是不会检测本地缓存是否过期了。

		强制使用本地缓存的漏洞：在redis不可用时，一个新的节点本地缓存它是获取不到数据的，可以考虑fail over来想其他节点拿去数据。

高可用也是可以通过本地缓存来实现的。
让本地缓存的过期时间 大于 redis的过期时间，当redis不可用的时候，那么本地缓存过期时间 - redis过期时间，这一段时间内起码服务是可用的。




本地缓存的过期时间、redis缓存的过期时间该如何设置?
	基本都是本地缓存过期时间 > redis缓存过期时间


多实例如何解决本地缓存问题？
	在实例下，只会有一个实例去计算榜单来更新缓存。

	对于已启动的实例：某个实例一旦计算完成榜单，可以采用Pub/Sub的模式让其他实例更新本地缓存。
	对于新启动的实例：新启动的实例本地缓存是不存在热榜数据的，因为可以采用fail over，在bff层请求某个实例，如果不存在则去其他实例获取。



*/

// 如何保证高可用?

// RankingService 定义榜单服务接口，除非你的榜单业务很复杂，那么可以抽象成单独的一个接口
type RankingService interface {
	// TopN 计算出N个排名
	TopN(ctx context.Context) error
}

//type compareFn[T any] func(src T, dst T)

// 我们所依赖的数据，可以通过repository获取，也可以聚合多个服务来获取。
type rankingService struct {
	repo repository.RankingRepository
	// 批量查询的大小
	batchSize int

	// 榜单长度
	length int

	interSvc InteractionService
	artSvc   ArticleService

	// 计算分数的函数
	scoreFn func(t time.Time, likes int64) float64
}

type node struct {
	Score   float64
	article domain.Article
}

func NewRankingService(
	artSvc ArticleService,
	interSvc InteractionService,
	repo repository.RankingRepository) RankingService {
	svc := &rankingService{
		repo:      repo,
		batchSize: 500,
		length:    10,
		artSvc:    artSvc,
		interSvc:  interSvc,
	}
	svc.scoreFn = svc.score
	return svc
}

func (svc *rankingService) score(utime time.Time, likes int64) float64 {
	duration := time.Since(utime).Seconds()
	return float64(likes-1) / math.Pow(duration+2, 1.5)
}

// TopN topn接口应该暴漏context，让外部来控制你执行的超市时间。
func (svc *rankingService) TopN(ctx context.Context) error {
	arts, err := svc.topN(ctx)
	if err != nil {
		return err
	}
	return svc.repo.ReplaceTopN(ctx, arts)
}

// TopN 文章的TopN计算
func (svc *rankingService) topN(ctx context.Context) ([]domain.Article, error) {
	// 为了让榜单的数据稳点一些，我们只计算7天内的数据。
	now := time.Now()
	offset := 0
	deadline := now.Add(-24 * 7 * time.Hour)

	container := queue.NewPriorityQueue[node](svc.length, func(src node, dst node) int {
		if src.Score > dst.Score {
			return 1
		} else if src.Score < dst.Score {
			return -1
		} else {
			return 0
		}
	})

	for {
		// now：now是为了秒顶记录，where time < ?，这样多个循环查询出来的记录总是相同的。
		articles, err := svc.artSvc.ListPub(ctx, now, offset, svc.batchSize)
		if err != nil {
			return nil, err
		}
		interIDs := slice.Map(articles, func(idx int, src domain.Article) int64 {
			return src.ID
		})
		inters, err := svc.interSvc.GetByIDs(ctx, "article", interIDs)
		if err != nil {
			return nil, err
		}

		nodes := slice.Map(articles, func(idx int, src domain.Article) node {
			return node{
				article: src,
				Score:   svc.scoreFn(src.UTime, inters[src.ID].Likes),
			}
		})

		for _, n := range nodes {
			if container.Len() < svc.length {
				_ = container.Enqueue(n)
				continue
			}
			// 取出一个元素，比较分数大小
			lastNode, _ := container.Dequeue()
			if n.Score > lastNode.Score {
				_ = container.Enqueue(n)
			} else {
				_ = container.Enqueue(lastNode)
			}
		}
		//查询出的数量不够一批的时候就结束循环
		if len(articles) < svc.batchSize {
			break
		}
		// 文章的更新时间在7天外也结束循环
		lastArt := articles[len(articles)-1]
		if lastArt.UTime.Before(deadline) {
			break
		}
		offset += svc.length
	}

	l := container.Len()
	result := make([]domain.Article, l)
	for i := l - 1; i >= 0; i-- {
		n, _ := container.Dequeue()
		result[i] = n.article
	}
	return result, nil
}
