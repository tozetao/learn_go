package web

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/service"
	svcmocks "learn_go/webook/internal/service/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		// 测试用例名
		name string

		// 模拟依赖的服务
		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService)

		// 构造的请求
		reqBuilder func(t *testing.T) *http.Request

		// 预期的结果
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "q123456",
				})
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{"email":"123@qq.com","password":"q123456"}`)))
				req.Header.Set("Content-Type", "application/json; charset=utf-8")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc, codeSvc := testCase.mock(ctrl)

			// 准备服务器，注册路由
			server := gin.Default()
			hld := NewUserHandler(userSvc, codeSvc)
			hld.RegisterRoutes(server)

			req := testCase.reqBuilder(t)
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			assert.Equal(t, testCase.wantCode, recorder.Code)
			assert.Equal(t, testCase.wantBody, recorder.Body.String())
		})
	}
}

func TestEmail(t *testing.T) {
	testCases := []struct {
		name   string
		email  string
		result bool
	}{
		{
			name:   "case1",
			email:  "test123456",
			result: false,
		},
		{
			name:   "case2",
			email:  "test123456@",
			result: false,
		},
		{
			name:   "后缀不完整1",
			email:  "test123456@qq",
			result: false,
		},
		{
			name:   "后缀不完整2",
			email:  "test123456@qq.",
			result: false,
		},
		{
			name:   "case3",
			email:  "test123@qq.com",
			result: true,
		},
		{
			name:   "邮箱名带-",
			email:  "test-123@qq.com",
			result: false,
		},
		{
			name:   "邮箱名带_和.",
			email:  "tes.t_123@qq.com",
			result: true,
		},
	}

	h := NewUserHandler(nil, nil, nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, _ := h.emailExp.MatchString(tc.email)
			// fmt.Printf("email: %s\n", tc.email)
			assert.Equal(t, tc.result, result)
		})
	}
}
