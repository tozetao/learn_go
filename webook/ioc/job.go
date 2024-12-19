package ioc

import (
	"github.com/robfig/cron/v3"
	"learn_go/webook/internal/job"
	"learn_go/webook/internal/service"
	"learn_go/webook/pkg/logger"
	"time"
)

func InitRankingJob(svc service.RankingService) *job.RankingJob {
	return job.NewRankingJob(svc, time.Second*30)
}

func InitCron(l logger.LoggerV2, rankingJob *job.RankingJob) *cron.Cron {
	// 任务超时时间 < 定时任务的间隔时间 < 缓存超时时间

	c := cron.New(cron.WithSeconds())
	builder := job.NewCronJobBuilder(l)

	// 每3分钟执行一次
	//_, err := c.AddJob("@every 3s", builder.Build(rankingJob))
	_, err := c.AddJob("0 */3 * * * ?", builder.Build(rankingJob))
	if err != nil {
		panic(err)
	}
	return c
}
