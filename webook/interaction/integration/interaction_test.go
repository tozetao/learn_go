package integration

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	intrv1 "learn_go/webook/api/proto/gen/intr"
	"learn_go/webook/interaction/integration/startup"
	"learn_go/webook/interaction/repository/dao"
	"strconv"
	"testing"
	"time"
)

/*
assert.NoError()

	断言一个错误值是否为nil。如果为nil，测试会继续执行，但会记录一个失败的测试结果

require.NoError()

	断言一个错误值是否为nil。如果为nil，测试会终止，不会执行后面的代码。
*/
var (
	likeCntField     = "like_cnt"
	readCntField     = "read_cnt"
	favoriteCntField = "favorite_cnt"
)

// 测试套件
type InteractionTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
	redis  redis.Cmdable
}

// hook, 在测试启动之前触发
func (s *InteractionTestSuite) SetupSuite() {
	s.db = startup.NewDB()

	s.redis = startup.NewRedis()

	//s.server = gin.Default()
	//s.server.Use(func(ctx *gin.Context) {
	//	ctx.Set("user", &web.UserClaims{
	//		Uid: 1000,
	//	})
	//})

	//// 新建handler
	//handler := startup.InitArticleHandler()
	//// 注册路由
	//handler.RegisterRoutes(s.server)

	t := s.T()
	err := s.db.Exec("truncate table `user_likes`").Error
	assert.NoError(t, err)
	err = s.db.Exec("truncate table `user_favorites`").Error
	assert.NoError(t, err)
	err = s.db.Exec("truncate table `interactions`").Error
	assert.NoError(t, err)

}

func (s *InteractionTestSuite) TearDownSuite() {
	//t := s.T()
	//
	//err := s.db.Exec("truncate table `articles`").Error
	//assert.NoError(t, err)
	//
	//err = s.db.Exec("truncate table `publish_articles`").Error
	//assert.NoError(t, err)
}

// 测试增加文章浏览量
func (s *InteractionTestSuite) TestView() {
	t := s.T()

	//authorID := 1001
	biz := "article"

	testCases := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		// 输入的数据
		// 业务类型
		biz string
		// 文章id
		bizID int64

		// 期望的输出
		wantErr  error
		wantResp *intrv1.ViewResp
	}{
		{
			name:  "db增加，缓存不增加",
			biz:   "article",
			bizID: 1,

			wantErr:  nil,
			wantResp: &intrv1.ViewResp{},

			before: func(t *testing.T) {
				err := s.db.Create(&dao.Interaction{
					ID:    1,
					Biz:   "article",
					BizID: 1,
					CTime: time.Now().UnixMilli()}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 查询文章的浏览量是否正常
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				// db有值
				bizID := 1
				var inter dao.Interaction
				err := s.db.WithContext(ctx).Model(&dao.Interaction{}).
					Where("biz = ? and biz_id = ?", biz, bizID).
					First(&inter).Error
				assert.NoError(t, err)
				assert.Equal(t, int64(1), inter.ReadCnt)

				// 缓存中空值
				key := s.key(biz, int64(bizID))
				m, _ := s.redis.HGetAll(ctx, key).Result()
				assert.Equal(t, "", m[likeCntField])
			},
		},
		{
			name:     "db, cache都有值",
			biz:      "article",
			bizID:    2,
			wantErr:  nil,
			wantResp: &intrv1.ViewResp{},

			before: func(t *testing.T) {
				var bizID int64 = 2
				err := s.db.Create(&dao.Interaction{
					ID:    2,
					Biz:   "article",
					BizID: bizID,
					CTime: time.Now().UnixMilli()}).Error
				assert.NoError(t, err)

				// 缓存中插入数据
				err = s.redis.HSet(context.Background(), s.key(biz, bizID),
					readCntField, 0).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var bizID int64 = 2
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				var inter dao.Interaction
				err := s.db.WithContext(ctx).Model(&dao.Interaction{}).
					Where("biz = ? and biz_id = ?", biz, bizID).
					First(&inter).Error
				assert.NoError(t, err)
				assert.Equal(t, int64(1), inter.ReadCnt)

				// 验证缓存中的值
				res, err := s.redis.HGet(ctx, s.key(biz, int64(bizID)), readCntField).Result()
				assert.NoError(t, err)
				readCnt, err := strconv.ParseInt(res, 10, 64)
				assert.NoError(t, err)
				assert.Equal(t, int64(1), readCnt)

				// 清空缓存
				err = s.redis.Del(context.Background(), s.key(biz, bizID)).Err()
				assert.NoError(t, err)
			},
		},
	}

	server := startup.InitInteractionServiceServer()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.before(t)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			// err := svc.View(ctx, testCase.biz, testCase.bizID)

			resp, err := server.View(ctx, &intrv1.ViewReq{
				Biz:   testCase.biz,
				BizId: testCase.bizID,
			})

			assert.NoError(t, err)
			assert.Equal(t, testCase.wantResp, resp)

			testCase.after(t)
		})
	}
}

func (s *InteractionTestSuite) key(biz string, bizID int64) string {
	return fmt.Sprintf("interaction:%s:%d", biz, bizID)
}

//// 测试点赞
//func (s *InteractionTestSuite) TestLike(t *testing.T) {
//
//}
//
//// 测试取消点赞
//func (s *InteractionTestSuite) TestUnlike(t *testing.T) {}

// 测试方法的运行入口
func TestInteraction(t *testing.T) {
	suite.Run(t, &InteractionTestSuite{})
}
