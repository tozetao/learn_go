package service

import (
	"context"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/repository"
	"learn_go/webook/pkg/logger"
	"time"
)

/*
抢占任务：
	抢占状态处于waiting的job，并将状态从waiting改为running

怎么知道一个节点是否正常?
	通过续约来判定一个节点是否仍然在执行中。
	续约就是更新任务的u_time，可以每个1分钟续约一次。如果一个节点处于running状态，且u_time超过1分钟没更新，那么该节点就失效了。

	p := now - 60s
	status = running and u_time < p

释放一个job
	将任务从running改为waiting

*/

type JobService interface {
	// Preempt 抢占一个任务
	Preempt(ctx context.Context) (domain.Job, error)

	Refresh(ctx context.Context, job domain.Job) error

	ResetNextTime(ctx context.Context, job domain.Job) error
}

func NewJobService(repo repository.JobRepository, l logger.LoggerV2) JobService {
	return &jobService{
		repo:     repo,
		l:        l,
		interval: time.Minute,
	}
}

func (svc *jobService) Preempt(ctx context.Context) (domain.Job, error) {
	j, err := svc.repo.Preempt(ctx)
	if err != nil {
		return domain.Job{}, err
	}

	// 续约
	timer := time.NewTicker(svc.interval)
	go func() {
		for range timer.C {
			ctx1, cancel := context.WithTimeout(context.Background(), time.Second)
			err1 := svc.Refresh(ctx1, j)
			cancel()
			if err1 != nil {
				svc.l.Error("job续约失败", logger.Int64("id", j.ID), logger.Error(err1))
			}
		}
	}()
	// 释放job
	j.CancelFunc = func() {
		timer.Stop()

		ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second)
		defer cancel2()

		err := svc.repo.Release(ctx2, j.ID)
		if err != nil {
			// 记录日志
			svc.l.Error("任务释放失败", logger.Int64("id", j.ID))
		}
	}
	return j, nil
}

func (svc *jobService) Refresh(ctx context.Context, job domain.Job) error {
	return svc.repo.UpdateUTime(ctx, job.ID, time.Now())
}

func (svc *jobService) ResetNextTime(ctx context.Context, job domain.Job) error {
	return svc.repo.UpdateNextTime(ctx, job.ID, job.NextTime())
}

type jobService struct {
	repo repository.JobRepository
	l    logger.LoggerV2

	// 间隔多久续约
	interval time.Duration
}
