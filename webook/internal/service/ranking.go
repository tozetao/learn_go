package service

import (
	"context"
	"github.com/ecodeclub/ekit/queue"
	"github.com/ecodeclub/ekit/slice"
	"learn_go/webook/internal/domain"
	"time"
)

// RankingService 定义榜单服务接口，除非你的榜单业务很复杂，那么可以抽象成单独的一个接口
type RankingService interface {
	// TopN 计算出N个排名
	TopN(n int) error
}

//type compareFn[T any] func(src T, dst T)

// 我们所依赖的数据，可以通过repository获取，也可以聚合多个服务来获取。
type rankingService struct {
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

func NewRankingService(batchSize int) RankingService {
	svc := &rankingService{
		batchSize: batchSize,
		length:    100,
	}
	return svc
}

func (svc *rankingService) TopN() error {
	return nil
}

// TopN 文章的TopN计算
func (svc *rankingService) topN() ([]domain.Article, error) {
	// 为了让榜单的数据稳点一些，我们只计算7天内的数据。
	now := time.Now()
	offset := 0

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
		articles, err := svc.artSvc.ListPub(context.Background(), now, offset, svc.batchSize)
		if err != nil {
			return nil, err
		}
		interIDs := slice.Map(articles, func(idx int, src domain.Article) int64 {
			return src.ID
		})
		inters, err := svc.interSvc.GetByIDs(context.Background(), "article", interIDs)
		if err != nil {
			return nil, err
		}

		nodes := slice.Map(articles, func(idx int, src domain.Article) node {
			return node{
				Score: svc.scoreFn(src.UTime, inters[src.ID].Likes),
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

		// 查询的记录包含7天外的数据就结束循环
		diffDays := now.Sub(articles[len(articles)-1].UTime).Hours() / 24
		if diffDays >= 7 || len(nodes) < svc.batchSize {
			break
		}
	}
	return nil, nil
}
