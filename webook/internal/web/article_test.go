package web

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/service"
	svcmocks "learn_go/webook/internal/service/mocks"
	"learn_go/webook/pkg/logger"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
单元测试：一个测试用例对应一个流程，更关注流程实现的细节。

 1. 定义测试模板
    确定有哪些输入、输出，有哪些要mock的服务。

 2. 实现测试用例的流程
    定义依赖服务接口
    mock出所依赖的服务接口
    实现测试用例要测试的流程代码。比如你测试article handler，你预期要使用哪些服务，调用哪些接口，这些都需要实现。gomock才能够mock出来。

3. 运行测试用例
*/
func TestArticleHandler_Publish(t *testing.T) {
	testCases := []struct {
		name string

		// 请求主体，json字符串。
		body string

		mock func(controller *gomock.Controller) service.ArticleService

		wantCode int
		wantRes  Result
	}{
		{
			name: "新建并发布",
			body: `
{
	"title": "This is my title",
	"content": "This is my content"
}
`,
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg:  "ok",
				Data: 1,
			},

			mock: func(controller *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(controller)

				// EXPECT: 返回一个对象，允许调用者指示预期用途。
				// Return: 声明了模拟函数返回的值。
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "This is my title",
					Content: "This is my content",
					Author: domain.Author{
						ID: 2001,
					},
				}).Return(int64(1), nil)

				return svc
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()
			server.Use(func(ctx *gin.Context) {
				ctx.Set("user", &UserClaims{
					Uid: 2001,
				})
			})

			articleSvc := testCase.mock(ctrl)
			articleHandler := NewArticleHandler(articleSvc, logger.NewNopLogger())
			articleHandler.RegisterRoutes(server)

			// 构建请求
			body := bytes.NewBufferString(testCase.body)
			req, err := http.NewRequest("POST", "/articles/publish", body)
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json; charset=utf-8")

			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			// 解析响应
			var res Result
			err = json.NewDecoder(resp.Body).Decode(&res)
			assert.NoError(t, err)

			// 断言响应结果
			assert.Equal(t, testCase.wantCode, resp.Code)
			if resp.Code != http.StatusOK {
				return
			}
			assert.Equal(t, testCase.wantRes, Result{
				Data: 1,
				Msg:  "ok",
			})
		})
	}

}
