package ratelimit

import (
	"context"
	"errors"
	"fmt"
	"learn_go/webook/internal/service/sms"
	"learn_go/webook/pkg/limiter/ratelimit"
)

var errLimited = errors.New("rate limit exceeded")

// RateLimitSMSService 针对整个Service做限流
type RateLimitSMSService struct {
	sms.Service

	limiter ratelimit.Limiter
	key     string
}

func NewRateLimitSMSService(svc sms.Service, l ratelimit.Limiter) *RateLimitSMSService {
	return &RateLimitSMSService{
		Service: svc,
		limiter: l,
		key:     "sms-limiter",
	}
}

func (r *RateLimitSMSService) Send(ctx context.Context, templateId string, params []string, phones []string) error {
	limited, err := r.limiter.Limit(ctx, r.key)

	// 发生错误是否限流? 从下游是否能够为你兜底进行考虑。下游强大，那么不限流；下游无法兜底则进行限流。
	if err != nil {
		return fmt.Errorf("系统错误, err: %v", err)
	}

	if limited {
		return errLimited
	}

	return r.Service.Send(ctx, templateId, params, phones)
}
