package job

import (
	"context"
	"errors"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/service"
	"learn_go/webook/pkg/logger"
	"time"
)

type Executor interface {
	Name() string

	Exec(ctx context.Context, j domain.Job) error
}

type LocalExecutor struct {
	funcs map[string]func(ctx context.Context, j domain.Job) error
}

func (l *LocalExecutor) Name() string {
	return "executor:local"
}

func (l *LocalExecutor) Exec(ctx context.Context, j domain.Job) error {
	fn, ok := l.funcs[j.Name]
	if !ok {
		return errors.New("找不到job的handler")
	}
	return fn(ctx, j)
}

func (l *LocalExecutor) RegisterFunc(name string, fn func(ctx context.Context, j domain.Job) error) {
	l.funcs[name] = fn
}

func NewLocalExecutor() *LocalExecutor {
	return &LocalExecutor{
		funcs: make(map[string]func(ctx context.Context, j domain.Job) error),
	}
}

type Scheduler struct {
	// job service
	svc service.JobService

	// 任务的最大执行时间
	jobDuration time.Duration

	// 每个多长时间进行续约
	interval time.Duration

	// 执行器列表
	executors map[string]Executor

	l logger.LoggerV2
}

func NewScheduler(svc service.JobService, l logger.LoggerV2) *Scheduler {
	return &Scheduler{
		jobDuration: time.Minute,
		interval:    time.Second * 30,
		executors:   make(map[string]Executor),
		svc:         svc,
		l:           l,
	}
}

func (s *Scheduler) Register(name string, executor Executor) {
	s.executors[name] = executor
}

// Schedule 调度任务, ctx决定了调度器什么时候结束
func (s *Scheduler) Schedule(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			break
		}

		// 1. 抢占一个job
		dbCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		j, err := s.svc.Preempt(dbCtx)
		cancel()
		if err != nil {
			continue
		}

		// 2. 找到该job所属的executor，执行该任务
		executor, ok := s.executors[j.Executor]
		if !ok {
			s.l.Error("找不到Executor", logger.Int64("job id", j.ID), logger.String("executor", j.Executor))
			continue
		}

		// 执行该任务
		go func() {
			execCtx, cancel := context.WithTimeout(context.Background(), s.jobDuration)
			defer cancel()

			err := executor.Exec(execCtx, j)
			if err != nil {
				s.l.Error("job执行失败", logger.Int64("job id", j.ID), logger.Error(err))

			}

			//// 释放job，重置job的执行时间
			//j.CancelFunc()
			//err = s.svc.ResetNextTime(, j)
			//if err != nil {
			//	s.l.Error("重置job执行时间失败", logger.Int64("job id", j.ID), logger.Error(err))
			//}
		}()

		// 这个循环对mysql压力太大了，我觉得应该要间隔的去抢占任务
		time.Sleep(time.Second)
	}
}
