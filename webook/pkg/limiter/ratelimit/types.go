package ratelimit

import "context"

type Limiter interface {
	// Limit 限流方法
	// key是要限流的对象（标识）
	// bool表示是否需要限流，true限流，false不限流。error表示限流器是否发生错误。
	Limit(ctx context.Context, key string) (bool, error)
}
