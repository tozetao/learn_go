package failover

import (
	"context"
	"learn_go/webook/internal/service/sms"
	"sync/atomic"
)

type TimeoutFailOverService struct {
	smsServices []sms.Service

	cnt int64
	idx int64

	// 阈值
	threshold int64
}

func NewTimeoutFailOverService(smsServices []sms.Service) *TimeoutFailOverService {
	return &TimeoutFailOverService{
		smsServices: smsServices,
		threshold:   5,
	}
}

// Send 非颜色的连续N个超时就切换服务，而是近似
func (t *TimeoutFailOverService) Send(ctx context.Context, templateId string, params []string, phones []string) error {
	idx := atomic.LoadInt64(&t.idx)
	cnt := atomic.LoadInt64(&t.cnt)

	// 这多个步骤不是原子操作，所有可能会有多个请求同时超过阈值
	if cnt > t.threshold {
		newIdx := idx % int64(len(t.smsServices))

		// CAS操作成功，说明有请求成功切换了。
		if atomic.CompareAndSwapInt64(&t.idx, idx, newIdx) {
			// 表示成功切换，计数要归0
			atomic.StoreInt64(&t.cnt, 0)
		}

		idx = atomic.LoadInt64(&t.idx)
	}

	err := t.smsServices[idx].Send(ctx, templateId, params, phones)
	switch err {
	case nil:
		atomic.StoreInt64(&t.idx, 0)
		return nil
	case context.DeadlineExceeded:
		// 超时
		atomic.AddInt64(&t.cnt, 1)
	default:
	}
	// 其他错误
	return err
}
