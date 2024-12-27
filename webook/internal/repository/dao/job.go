package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

const (
	jobStatusUnknown int = iota
	jobStatusWaiting
	jobStatusRunning
	jobStatusPause
)

type Job struct {
	ID int64 `json:"id" gorm:"primaryKey, autoincrement"`

	Name string `json:"name" gorm:"unique_index"`

	Cfg string `json:"cfg"`

	// Job所属执行器的名字
	Executor string
	// cron表达式
	Expression string

	Status int8 `json:"status" gorm:"column:status;index:status_next_time"`
	// 下一次执行的时间
	NextTime int64 `json:"next_time" gorm:"column:next_time;index:status_next_time"`

	Version string

	CTime int64 `json:"c_time" gorm:"column:c_time"`
	UTime int64 `json:"u_time" gorm:"column:u_time"`
}

type JobDao interface {
	Preempt(ctx context.Context) (Job, error)

	UpdateUTime(ctx context.Context, id int64, now time.Time) error
	UpdateNextTime(ctx context.Context, id int64, nt time.Time) error

	Release(ctx context.Context, id int64) error
}

type jobDao struct {
	db *gorm.DB
}

func (dao *jobDao) Release(ctx context.Context, id int64) error {
	return dao.db.WithContext(ctx).Model(&Job{}).Where("id=?", id).
		Where("id = ? and status = ?", id, jobStatusRunning).
		Update("status", jobStatusWaiting).Error
}

func (dao *jobDao) UpdateUTime(ctx context.Context, id int64, now time.Time) error {
	return dao.db.WithContext(ctx).Model(&Job{}).Where("id=?", id).Updates(map[string]interface{}{
		"u_time": now.UnixMilli(),
	}).Error
}

func (dao *jobDao) UpdateNextTime(ctx context.Context, id int64, nt time.Time) error {
	return dao.db.WithContext(ctx).Model(&Job{}).Where("id=?", id).Updates(map[string]interface{}{
		"next_time": nt.UnixMilli(),
	}).Error
}

func (dao *jobDao) Preempt(ctx context.Context) (Job, error) {
	now := time.Now()
	var job Job
	for {
		err := dao.db.WithContext(ctx).Model(&Job{}).
			Where("status = ? and next_time < ?", jobStatusWaiting, now.UnixMilli()).
			First(&job).Error
		// 发生错误 或者 找不到job就返回
		if err != nil {
			return Job{}, err
		}
		// 利用锁，将job从waiting变为running状态，保证只有一个实例抢占到该任务。
		res := dao.db.WithContext(ctx).Model(&Job{}).
			Where("id = ? and status = ?", job.ID, jobStatusWaiting).
			Updates(map[string]interface{}{
				"status": jobStatusRunning,
				"u_time": now.UnixMilli(),
			})
		if res.Error != nil {
			return Job{}, res.Error
		}
		if res.RowsAffected != 1 {
			continue
		}
		return job, nil
	}
}

func NewJobDao(db *gorm.DB) JobDao {
	return &jobDao{
		db: db,
	}
}
