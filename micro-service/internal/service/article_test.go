package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/repository/article"
	artrepomocks "learn_go/webook/internal/repository/mocks/article"
	"learn_go/webook/pkg/logger"
	"testing"
)

// 一个测试用例对应着代码处理流程的一个分支。

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
			article: domain.Article{
				Title:   "Title",
				Content: "Content",
				Author: domain.Author{
					ID: 2000,
				},
			},
			wantErr: nil,
			wantId:  5,
			mock: func(controller *gomock.Controller) (article.AuthorRepository, article.ReaderRepository) {
				author := artrepomocks.NewMockAuthorRepository(controller)
				reader := artrepomocks.NewMockReaderRepository(controller)

				art := domain.Article{
					Title:   "Title",
					Content: "Content",
					Author: domain.Author{
						ID: 2000,
					},
				}
				author.EXPECT().Create(gomock.Any(), art).Return(int64(5), nil)

				art.ID = 5
				reader.EXPECT().Save(gomock.Any(), art).Return(int64(5), nil)

				return author, reader
			},
		},
		{
			name: "新建失败",
			article: domain.Article{
				Title:   "Title",
				Content: "Content",
				Author: domain.Author{
					ID: 2000,
				},
			},
			wantErr: errors.New("mock db error"),
			wantId:  0,
			mock: func(controller *gomock.Controller) (article.AuthorRepository, article.ReaderRepository) {
				author := artrepomocks.NewMockAuthorRepository(controller)
				reader := artrepomocks.NewMockReaderRepository(controller)

				art := domain.Article{
					Title:   "Title",
					Content: "Content",
					Author: domain.Author{
						ID: 2000,
					},
				}
				author.EXPECT().Create(gomock.Any(), art).Return(int64(0), errors.New("mock db error"))
				return author, reader
			},
		},
		{
			name: "新建成功，发布失败后重试成功",
			article: domain.Article{
				Title:   "Title",
				Content: "Content",
				Author: domain.Author{
					ID: 2000,
				},
			},
			wantErr: nil,
			wantId:  5,
			mock: func(controller *gomock.Controller) (article.AuthorRepository, article.ReaderRepository) {
				author := artrepomocks.NewMockAuthorRepository(controller)
				reader := artrepomocks.NewMockReaderRepository(controller)

				art := domain.Article{
					Title:   "Title",
					Content: "Content",
					Author: domain.Author{
						ID: 2000,
					},
				}
				author.EXPECT().Create(gomock.Any(), art).Return(int64(5), nil)

				art.ID = 5
				reader.EXPECT().Save(gomock.Any(), art).Return(int64(0), errors.New("mock db error"))

				reader.EXPECT().Save(gomock.Any(), art).Return(int64(5), nil)

				return author, reader
			},
		},
		{
			name: "新建成功，发布失败，且次数用尽",
			article: domain.Article{
				Title:   "Title",
				Content: "Content",
				Author: domain.Author{
					ID: 2000,
				},
			},
			wantErr: errors.New("failed to publish"),
			wantId:  0,
			mock: func(controller *gomock.Controller) (article.AuthorRepository, article.ReaderRepository) {
				author := artrepomocks.NewMockAuthorRepository(controller)
				reader := artrepomocks.NewMockReaderRepository(controller)

				art := domain.Article{
					Title:   "Title",
					Content: "Content",
					Author: domain.Author{
						ID: 2000,
					},
				}
				author.EXPECT().Create(gomock.Any(), art).Return(int64(5), nil)

				art.ID = 5
				reader.EXPECT().Save(gomock.Any(), art).Times(3).
					Return(int64(0), errors.New("mock db error"))

				return author, reader
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authorRepo, readerRepo := testCase.mock(ctrl)
			svc := NewArticleService(nil, authorRepo, readerRepo, nil, logger.NewNopLogger())
			id, err := svc.PublishV1(context.Background(), testCase.article)

			assert.Equal(t, testCase.wantErr, err)
			assert.Equal(t, testCase.wantId, id)
		})
	}

}
