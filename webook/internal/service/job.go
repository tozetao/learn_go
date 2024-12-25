package service

import "context"

type JobService interface {

	// Preempt 抢占一个任务
	Preempt(ctx context.Context)
}

type jobService struct {
}
