package job

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"learn_go/webook/pkg/logger"
	"strconv"
	"time"
)

// Job 任务接口
type Job interface {
	Name() string
	Run() error
}

// 构建cron对象，同时利用适配器，监控任务的性能

type jobAdapterFunc func()

func (fn jobAdapterFunc) Run() {
	fn()
}

type CronJob struct {
	l   logger.LoggerV2
	vec *prometheus.SummaryVec
}

// Build 构建cron job对象
func (c *CronJob) Build(job Job) cron.Job {
	return jobAdapterFunc(func() {
		now := time.Now()

		err := job.Run()

		c.vec.WithLabelValues(job.Name(), strconv.FormatBool(err == nil)).
			Observe(float64(time.Since(now).Milliseconds()))
	})
}

func NewCronJob(l logger.LoggerV2) *CronJob {
	return &CronJob{
		l: l,
		vec: prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace: "go_project",
			Subsystem: "webook",
			Name:      "cron_job",
			Help:      "定时任务的性能检测",
			Objectives: map[float64]float64{
				0.5:  0.01,
				0.95: 0.01,
				0.99: 0.005,
			},
		}, []string{"job_name", "success"}),
	}
}
