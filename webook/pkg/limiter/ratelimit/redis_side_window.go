package ratelimit

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed slide_window.lua
var script string

type RedisSideWindow struct {
	cmd redis.Cmdable

	// 窗口大小
	interval time.Duration

	// 速率
	rate int

	// 表示在指定的interval时间内（窗口内）允许有rate个请求。
}

func NewRedisSideWindow(cmd redis.Cmdable, interval time.Duration, rate int) *RedisSideWindow {
	return &RedisSideWindow{
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}

func (r *RedisSideWindow) Limit(ctx context.Context, key string) (bool, error) {
	return r.cmd.Eval(ctx, script, []string{key},
		r.interval.Milliseconds(), r.rate, time.Now().UnixMilli()).Bool()
}
