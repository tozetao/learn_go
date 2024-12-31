package service

import (
	"context"
	"github.com/ecodeclub/ekit/queue"
	"github.com/ecodeclub/ekit/slice"
	service2 "learn_go/webook/interaction/service"
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

	redis、mysql都不可用：
		利用本地缓存来做兜底。
		Get接口的逻辑是，先从本地缓存获取，再去redis中获取。
		我们可以这样修改Get接口逻辑，当redis不可用时，再次从本地缓存中获取，此时就不监测本地缓存的过期时间了。

		强制使用本地缓存的漏洞：在redis不可用时，一个新的节点本地缓存它是获取不到数据的，可以考虑fail over来想其他节点拿去数据。


本地缓存的过期时间、redis缓存的过期时间该如何设置?
	个人认为只要俩个过期时间不同就可以了。


多实例如何解决本地缓存问题？
	在redis正常情况下，每个实例判定本地缓存不存在就会去redis获取，这是正常逻辑。
	而当redis不可用时，旧的实例可以返回本地缓存中过期的数据，新的实例只能使用fail over去其他实例获取其他数据了。


*/

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

	interSvc service2.InteractionService
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
	interSvc service2.InteractionService,
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
