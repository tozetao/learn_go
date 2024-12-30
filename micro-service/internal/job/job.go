package job

import (
	"context"
	"learn_go/webook/internal/service"
	"time"
)

// 1. 先实现自己定义的job接口。2.将自定义job转成cron job 3. 创建cron对象，加入所有创建的job

// Job 任务接口
type Job interface {
	Name() string
	Run() error
}

// RankingJob 排行榜任务
type RankingJob struct {
	rankingSvc service.RankingService
	duration   time.Duration
}

func NewRankingJob(rankingSvc service.RankingService, duration time.Duration) *RankingJob {
	return &RankingJob{rankingSvc: rankingSvc, duration: duration}
}

func (r *RankingJob) Name() string {
	return "ranking"
}

func (r *RankingJob) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), r.duration)
	defer cancel()
	return r.rankingSvc.TopN(ctx)
}
