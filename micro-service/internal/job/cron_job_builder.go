package job

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"learn_go/webook/pkg/logger"
	"strconv"
	"time"
)

type jobAdapterFunc func()

func (fn jobAdapterFunc) Run() {
	fn()
}

type CronJobBuilder struct {
	l   logger.LoggerV2
	vec *prometheus.SummaryVec
}

// Build 将Job任务转换成cron.Job
func (c *CronJobBuilder) Build(job Job) cron.Job {
	return jobAdapterFunc(func() {
		c.l.Info(fmt.Sprintf("job:%s execution in progress", job.Name()))
		var success bool
		now := time.Now()

		defer func() {
			since := time.Since(now).Milliseconds()

			c.vec.WithLabelValues(job.Name(), strconv.FormatBool(success)).
				Observe(float64(since))
		}()

		err := job.Run()
		if err != nil {
			c.l.Error(fmt.Sprintf("job:%s run error", job.Name()), logger.Error(err))
		} else {
			c.l.Info(fmt.Sprintf("job:%s execution completed", job.Name()))
		}
		success = err == nil
	})
}

func NewCronJobBuilder(l logger.LoggerV2) *CronJobBuilder {
	return &CronJobBuilder{
		l: l,
		vec: prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace: "go_project",
			Subsystem: "webook",
			Name:      "jobs",
			Help:      "定时任务的性能测量",
			Objectives: map[float64]float64{
				0.5:  0.01,
				0.95: 0.01,
				0.99: 0.005,
			},
		}, []string{"job_name", "success"}),
	}
}
