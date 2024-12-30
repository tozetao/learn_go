package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"learn_go/webook/internal/repository/cache/redismocks"
	"testing"
)

func TestRedisCodeCache_Set(t *testing.T) {
	biz := "login"
	phone := "13512341234"
	code := "123456"
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) redis.Cmdable

		ctx   context.Context
		biz   string
		phone string
		code  string

		wantErr error
	}{
		{
			name: "设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)

				// 这是Eval的返回值
				redisCmd := redis.NewCmd(context.Background())
				redisCmd.SetErr(nil)
				redisCmd.SetVal(int64(0))

				key := fmt.Sprintf("%s:code:%s", biz, phone)
				cmd.EXPECT().Eval(context.Background(), sendCodeScript, []string{key}, []any{code, 60 * 10}).Return(redisCmd)

				return cmd
			},
			ctx:   context.Background(),
			biz:   biz,
			phone: phone,
			code:  code,
		},
		{
			name: "未知错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)

				// 这是Eval的返回值
				redisCmd := redis.NewCmd(context.Background())
				redisCmd.SetErr(errors.New("unknown error"))
				redisCmd.SetVal(int64(-2))

				key := fmt.Sprintf("%s:code:%s", biz, phone)
				cmd.EXPECT().Eval(context.Background(), sendCodeScript, []string{key}, []any{code, 60 * 10}).Return(redisCmd)

				return cmd
			},
			ctx:   context.Background(),
			biz:   biz,
			phone: phone,
			code:  code,

			wantErr: errors.New("unknown error"),
		},
		// 缓存错误
		// 发送太多次错误
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			codeCache := NewCodeCache(tc.mock(ctrl))

			err := codeCache.Set(tc.ctx, tc.biz, tc.phone, tc.code)

			assert.Equal(t, tc.wantErr, err)
		})
	}
}
