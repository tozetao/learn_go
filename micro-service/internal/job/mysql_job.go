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
			s.l.Info("抢占不到任务", logger.Error(err))
			time.Sleep(time.Second)
			continue
		}

		// 2. 找到该job所属的executor，执行该任务
		executor, ok := s.executors[j.Executor]
		if !ok {
			s.l.Error("找不到Executor", logger.Int64("job id", j.ID), logger.String("executor", j.Executor))
			time.Sleep(time.Second)
			continue
		}
		s.l.Info("开始执行任务", logger.Int64("job id", j.ID), logger.String("next_time", j.Nt.Format(time.DateTime)))

		// 执行该任务
		// TODO: 外部不知道任务是否执行完毕。因为目前开启一个goroutine来执行job，通过context设置goroutine的最大执行时间。
		go func() {
			execCtx, cancel := context.WithTimeout(context.Background(), s.jobDuration)
			defer func() {
				s.l.Info("任务执行完毕")
				cancel()
				j.CancelFunc()
			}()

			err := executor.Exec(execCtx, j)
			if err != nil {
				s.l.Error("job执行失败", logger.Int64("job id", j.ID), logger.Error(err))
			}

			// 重置job的执行时间
			err = s.svc.ResetNextTime(execCtx, j)
			if err != nil {
				s.l.Error("重置job执行时间失败", logger.Int64("job id", j.ID), logger.Error(err))
			}
		}()

		time.Sleep(time.Millisecond * 500)
	}

}
