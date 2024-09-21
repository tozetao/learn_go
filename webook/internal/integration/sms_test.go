package integration

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"learn_go/webook/ioc"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSMSHandler_Send(t *testing.T) {
	// 初始化服务器
	server := InitWebServer("test_template")
	rdb := ioc.NewRedis()

	testCases := []struct {
		name string

		// 传入的参数
		phone      string
		reqBuilder func(t *testing.T) *http.Request
		// 准备之前
		before func(t *testing.T)

		// 准备之后
		after func(t *testing.T)

		wantCode   int
		wantResult string
	}{
		{
			// 发送成功
			name: "发送成功",

			phone: "13512341234",
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/sms/send", bytes.NewReader([]byte(`{"phone": "13512341234"}`)))
				req.Header.Set("Content-Type", "application/json; charset=utf-8")
				assert.NoError(t, err)
				return req
			},

			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				// 清空发送的验证码
				key := fmt.Sprintf("%s:code:%s", "login", "13512341234")
				code, err := rdb.Get(ctx, key).Result()
				// t.Logf("code: %v, err:%v\n", code, err)
				assert.NoError(t, err)
				assert.True(t, len(code) == 6)

				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},

			wantCode:   http.StatusOK,
			wantResult: "success",
		},
		{
			name: "发送太多次了",

			phone: "13512341234",
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/sms/send", bytes.NewReader([]byte(`{"phone": "13512341234"}`)))
				req.Header.Set("Content-Type", "application/json; charset=utf-8")
				assert.NoError(t, err)
				return req
			},

			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				// 预设一个key进行占位
				key := fmt.Sprintf("%s:code:%s", "login", "13512341234")
				err := rdb.Set(ctx, key, "123456", time.Minute*9+time.Second*30).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				key := fmt.Sprintf("%s:code:%s", "login", "13512341234")

				code, err := rdb.Get(ctx, key).Result()
				assert.Equal(t, code, "123456")
				assert.NoError(t, err)

				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},

			wantCode:   http.StatusOK,
			wantResult: "验证码发送太多次了",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)

			// 构建请求
			req := tc.reqBuilder(t)

			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantResult, recorder.Body.String())

			tc.after(t)
		})
	}
}
