package service

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"learn_go/webook/internal/domain"
	svcmocks "learn_go/webook/internal/service/mocks"
	"testing"
	"time"
)

// 单元测试
func TestRankingService(t *testing.T) {
	now := time.Now()

	utime := now.Add(time.Hour * -12)

	testCases := []struct {
		Name string

		// 输入

		// mock的服务
		mock func(controller *gomock.Controller) (ArticleService, InteractionService)

		// 期待的输出
		wantErr    error
		wantResult []domain.Article
	}{
		{
			Name: "测试榜单名次",
			mock: func(controller *gomock.Controller) (ArticleService, InteractionService) {
				artSvc := svcmocks.NewMockArticleService(controller)
				interSvc := svcmocks.NewMockInteractionService(controller)

				artSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]domain.Article{
						{ID: 1, UTime: utime},
						{ID: 2, UTime: utime},
						{ID: 3, UTime: utime},
					}, nil)

				interSvc.EXPECT().GetByIDs(gomock.Any(), "article", []int64{1, 2, 3}).
					Return(map[int64]domain.Interaction{
						1: {ID: 101, Biz: "article", BizID: 1, Likes: 1},
						2: {ID: 102, Biz: "article", BizID: 2, Likes: 2},
						3: {ID: 103, Biz: "article", BizID: 3, Likes: 3},
					}, nil)
				return artSvc, interSvc
			},
			wantErr: nil,
			wantResult: []domain.Article{
				{ID: 3, UTime: utime},
				{ID: 2, UTime: utime},
				{ID: 1, UTime: utime},
			},
		},
	}

	for _, testCase := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		artSvc, interSvc := testCase.mock(ctrl)

		rankingSvc := NewRankingService(artSvc, interSvc)
		result, err := rankingSvc.topN()
		assert.Equal(t, testCase.wantErr, err)
		assert.Equal(t, testCase.wantResult, result)
	}
}
