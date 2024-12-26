package domain

import (
	"github.com/robfig/cron/v3"
	"time"
)

/*

Job
	定时任务的抽象。

Executor
	执行器，注册处理job的handler

Scheduler
	调度器。负责取出任务，并执行任务。

*/

type Job struct {
	ID   int64
	Name string
	Cfg  string

	// Job所属执行器的名字
	Executor string

	// cron表达式
	Expression string

	CancelFunc func()
}

// NextTime 返回下一次的执行时间
func (j Job) NextTime() time.Time {
	parser := cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)
	s, _ := parser.Parse(j.Expression)
	return s.Next(time.Now())
}
