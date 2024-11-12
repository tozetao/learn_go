package integration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"learn_go/webook/internal/integration/startup"
	"learn_go/webook/internal/repository/dao"
	"learn_go/webook/internal/web"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
assert.NoError()
	断言一个错误值是否为nil。如果为nil，测试会继续执行，但会记录一个失败的测试结果
require.NoError()
	断言一个错误值是否为nil。如果为nil，测试会终止，不会执行后面的代码。
*/

type Article struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	//// 作者
	AuthorID int64 `json:"author_id"`
	Ctime    int64 `json:"c_time"`
	Utime    int64 `json:"u_time"`
}

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

// 测试套件
type ArticleTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

// hook, 在测试启动之前触发
func (s *ArticleTestSuite) SetupSuite() {
	s.db = startup.NewDB()

	s.server = gin.Default()
	s.server.Use(func(ctx *gin.Context) {
		ctx.Set("user", &web.UserClaims{
			Uid: 1000,
		})
	})

	// 新建handler
	handler := startup.InitArticleHandler()
	// 注册路由
	handler.RegisterRoutes(s.server)
}

func (s *ArticleTestSuite) TearDownSuite() {
	t := s.T()

	err := s.db.Exec("truncate table `articles`").Error
	assert.NoError(t, err)
}

// 定义测试方法
func (s *ArticleTestSuite) TestFoo() {
	s.T().Log("hello, 这是测试套件")
}

func (s *ArticleTestSuite) TestEdit() {
	t := s.T()

	testCases := []struct {
		name string

		reqBuilder func(t *testing.T, article Article) *http.Request
		before     func(t *testing.T)
		after      func(t *testing.T)

		// 输入的数据
		article Article

		// 期望的输出
		wantCode int
		wantRes  Result[int64]
	}{
		{
			name: "新建帖子",
			reqBuilder: func(t *testing.T, article Article) *http.Request {
				buf, err := json.Marshal(article)
				assert.NoError(t, err)

				// 构建起请求
				req, err := http.NewRequest("POST", "/articles/edit", bytes.NewBuffer(buf))
				assert.NoError(t, err)

				req.Header.Set("Content-Type", "application/json; charset=utf-8")
				return req
			},
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				// 核对数据

				// 查询ID=1的记录
				var article dao.Article
				var authorID int64 = 1000
				var articleID int64 = 1

				err := s.db.Where("id", 1).First(&article).Error
				assert.NoError(t, err)

				assert.Equal(t, articleID, article.ID)
				assert.Equal(t, "hello", article.Title)
				assert.Equal(t, authorID, article.AuthorID)
				assert.True(t, article.Ctime > 0)
				assert.True(t, article.Utime > 0)
			},

			article: Article{
				Title:   "hello",
				Content: "This is content.",
			},

			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "ok",
				Data: 1,
			},
		},
		{
			name: "修改帖子",
			reqBuilder: func(t *testing.T, article Article) *http.Request {
				buf, err := json.Marshal(article)
				assert.NoError(t, err)

				// 构建起请求
				req, err := http.NewRequest("POST", "/articles/edit", bytes.NewBuffer(buf))
				assert.NoError(t, err)

				req.Header.Set("Content-Type", "application/json; charset=utf-8")
				return req
			},

			article: Article{
				ID:      3,
				Title:   "new article",
				Content: "new content",
			},

			before: func(t *testing.T) {
				// 要准备一个已经存在的帖子
				article := dao.Article{
					ID:       3,
					Title:    "my article",
					Content:  "my content",
					AuthorID: 1000,
					Ctime:    123,
					Utime:    234,
				}
				err := s.db.Create(&article).Error
				assert.NoError(t, err)
			},

			after: func(t *testing.T) {
				// 查询ID=1的记录
				var article dao.Article
				var authorID int64 = 1000
				var articleID int64 = 3

				err := s.db.Where("id", articleID).First(&article).Error
				assert.NoError(t, err)

				// 更新后的时间必定会大于准备的Utime
				assert.True(t, article.Utime > 234)
				article.Utime = 0

				assert.Equal(t, dao.Article{
					ID:       articleID,
					Title:    "new article",
					Content:  "new content",
					Ctime:    123,
					AuthorID: authorID,
				}, article)
			},

			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "ok",
				Data: 3,
			},
		},
		{
			name: "修改帖子 - 别人的帖子",
			reqBuilder: func(t *testing.T, article Article) *http.Request {
				buf, err := json.Marshal(article)
				assert.NoError(t, err)

				// 构建起请求
				req, err := http.NewRequest("POST", "/articles/edit", bytes.NewBuffer(buf))
				assert.NoError(t, err)

				req.Header.Set("Content-Type", "application/json; charset=utf-8")
				return req
			},

			article: Article{
				ID:      15,
				Title:   "new article",
				Content: "new content",
			},

			before: func(t *testing.T) {
				// 插入一个别人的帖子
				article := dao.Article{
					ID:       15,
					Title:    "my article",
					Content:  "my content",
					AuthorID: 850,
					Ctime:    123,
					Utime:    234,
				}
				err := s.db.Create(&article).Error
				assert.NoError(t, err)
			},

			after: func(t *testing.T) {
				// 由于是别人的帖子，修改应该是失败的，因此id=3的帖子的数据是不变的。

				var article dao.Article
				var authorID int64 = 850
				var articleID int64 = 15

				err := s.db.Where("id", articleID).First(&article).Error
				assert.NoError(t, err)

				assert.Equal(t, dao.Article{
					ID:       articleID,
					Title:    "my article",
					Content:  "my content",
					Ctime:    123,
					Utime:    234,
					AuthorID: authorID,
				}, article)
			},

			wantCode: http.StatusInternalServerError,
			wantRes: Result[int64]{
				Msg:  "failed",
				Code: 5,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.before(t)

			// 1. 构建请求
			req := testCase.reqBuilder(t, testCase.article)

			// 2. 处理该请求并写入响应
			resp := httptest.NewRecorder()
			s.server.ServeHTTP(resp, req)

			// 3. 校验响应
			assert.Equal(t, testCase.wantCode, resp.Code)

			var webResult Result[int64]
			err := json.Unmarshal(resp.Body.Bytes(), &webResult)
			require.NoError(t, err)

			assert.Equal(t, testCase.wantRes.Msg, webResult.Msg)
			assert.Equal(t, testCase.wantRes.Code, webResult.Code)
			assert.Equal(t, testCase.wantRes.Data, webResult.Data)

			testCase.after(t)
		})
	}
}

// 测试方法的运行入口
func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
}
