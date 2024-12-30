package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/send_code.lua
	sendCodeScript string

	//go:embed lua/verify_code.lua
	verifyCodeScript string

	ErrTooManyVerify = errors.New("code verify too many times")
	ErrTooManySend   = errors.New("code send too many times")
)

type CodeCache interface {
	Set(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

// RedisCodeCache 验证码缓存
type RedisCodeCache struct {
	cache redis.Cmdable
}

func NewCodeCache(cache redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		cache: cache,
	}
}

func (c *RedisCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	result, err := c.cache.Eval(ctx, sendCodeScript, []string{c.key(biz, phone)}, code, 10*60).Int()
	// 缓存错误
	if err != nil {
		return err
	}

	switch result {
	case -2:
		return errors.New("unknown error")
	case -1:
		return ErrTooManySend
	}
	return nil
}

func (c *RedisCodeCache) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	result, err := c.cache.Eval(ctx, verifyCodeScript, []string{c.key(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch result {
	case -2:
		return false, ErrTooManyVerify
	case -1:
		return false, nil
	}
	return true, nil
}

func (c *RedisCodeCache) key(biz string, phone string) string {
	return fmt.Sprintf("%s:code:%s", biz, phone)
}
