package service

import (
	"go.uber.org/mock/gomock"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/repository/article"
	artrepomocks "learn_go/webook/internal/repository/mocks/article"
	"testing"
)

func Test_articleService_PublishV1(t *testing.T) {
	testCases := []struct {
		name string

		// 输入
		article domain.Article

		// 输出
		wantErr error
		wantId  int64

		mock func(controller *gomock.Controller) (article.AuthorRepository, article.ReaderRepository)
	}{
		{
			name: "新建并发布",
			
			mock: func(controller *gomock.Controller) (article.AuthorRepository, article.ReaderRepository) {
				author := artrepomocks.NewMockAuthorRepository(controller)
				reader := artrepomocks.NewMockReaderRepository(controller)

				author.EXPECT().Create(gomock.Any(), domain.Article{
					Title:   "Title",
					Content: "Content",
					Author: domain.Author{
						ID: 2000,
					},
				}).Return(int64(5), nil)

				reader.EXPECT().Save(gomock.Any(), domain.Article{
					ID:      5,
					Title:   "Title",
					Content: "Content",
					Author: domain.Author{
						ID: 2000,
					},
				}).Return(int64(5), nil)

				return author, reader
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			//articleSvc := testCase.mock(ctrl)
			//articleHandler := NewArticleHandler(articleSvc, logger.NewNopLogger())
			//articleHandler.RegisterRoutes(server)
			//
			//// 断言响应结果
			//assert.Equal(t, testCase.wantCode, resp.Code)
			//if resp.Code != http.StatusOK {
			//	return
			//}
			//assert.Equal(t, testCase.wantRes, Result{
			//	Data: 1,
			//	Msg:  "ok",
			//})
		})
	}

}
