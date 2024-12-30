package repository

import (
	"context"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/repository/dao"
	"time"
)

type JobRepository interface {
	// Preempt 抢占一个任务
	Preempt(ctx context.Context) (domain.Job, error)

	Release(ctx context.Context, id int64) error

	UpdateUTime(ctx context.Context, id int64, now time.Time) error

	UpdateNextTime(ctx context.Context, id int64, nt time.Time) error
}

type CronJobRepository struct {
	dao dao.JobDao
}

func (repo *CronJobRepository) UpdateUTime(ctx context.Context, id int64, now time.Time) error {
	return repo.dao.UpdateUTime(ctx, id, now)
}

func (repo *CronJobRepository) UpdateNextTime(ctx context.Context, id int64, nt time.Time) error {
	return repo.dao.UpdateNextTime(ctx, id, nt)
}

func (repo *CronJobRepository) Preempt(ctx context.Context) (domain.Job, error) {
	j, err := repo.dao.Preempt(ctx)
	if err != nil {
		return domain.Job{}, err
	}
	return domain.Job{
		ID:         j.ID,
		Name:       j.Name,
		Cfg:        j.Cfg,
		Executor:   j.Executor,
		Expression: j.Expression,
		Nt:         time.UnixMilli(j.NextTime),
	}, err
}

func (repo *CronJobRepository) Release(ctx context.Context, id int64) error {
	return repo.dao.Release(ctx, id)
}

func NewCronJobRepository(dao dao.JobDao) JobRepository {
	return &CronJobRepository{
		dao: dao,
	}
}
