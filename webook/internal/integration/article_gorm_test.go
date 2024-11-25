package integration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/integration/startup"
	"learn_go/webook/internal/repository/dao"
	"learn_go/webook/internal/web"
	"learn_go/webook/pkg/ginx"
	"net/http"
	"net/http/httptest"
	"net/url"
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

	t := s.T()
	err := s.db.Exec("truncate table `articles`").Error
	assert.NoError(t, err)
	err = s.db.Exec("truncate table `publish_articles`").Error
	assert.NoError(t, err)
}

func (s *ArticleTestSuite) TearDownSuite() {
	//t := s.T()
	//
	//err := s.db.Exec("truncate table `articles`").Error
	//assert.NoError(t, err)
	//
	//err = s.db.Exec("truncate table `publish_articles`").Error
	//assert.NoError(t, err)
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
			article: Article{
				Title:   "hello",
				Content: "This is content.",
			},

			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "ok",
				Data: 1,
			},
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

				err := s.db.First(&article).Error
				assert.NoError(t, err)

				assert.True(t, article.ID > 0)
				assert.True(t, article.Ctime > 0)
				assert.True(t, article.Utime > 0)

				article.ID = 0
				article.Ctime = 0
				article.Utime = 0

				assert.Equal(t, dao.Article{
					Title:    "hello",
					Content:  "This is content.",
					AuthorID: 1000,
					Status:   domain.ArticleStatusUnpublished,
				}, article)
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
				ID:      5,
				Title:   "new title",
				Content: "new content",
			},

			before: func(t *testing.T) {
				// 要准备一个已经存在的帖子
				article := dao.Article{
					ID:       5,
					Title:    "title",
					Content:  "content",
					Status:   domain.ArticleStatusUnpublished,
					AuthorID: 1000,
					Ctime:    123,
					Utime:    234,
				}
				err := s.db.Create(&article).Error
				assert.NoError(t, err)
			},

			after: func(t *testing.T) {
				var article dao.Article

				err := s.db.Where("id", 5).First(&article).Error
				assert.NoError(t, err)

				// 更新后的时间必定会大于准备的Utime
				assert.True(t, article.Ctime == 123)
				assert.True(t, article.Utime > 234)
				article.Utime = 0

				assert.Equal(t, dao.Article{
					ID:       5,
					Title:    "new title",
					Content:  "new content",
					Status:   domain.ArticleStatusUnpublished,
					Ctime:    123,
					AuthorID: 1000,
				}, article)
			},

			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "ok",
				Data: 5,
			},
		},
		{
			name: "修改别人的帖子",
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
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "failed",
				Code: 5,
			},

			before: func(t *testing.T) {
				// 插入一个别人的帖子
				article := dao.Article{
					ID:       15,
					Title:    "my article",
					Content:  "my content",
					Status:   domain.ArticleStatusPublished,
					AuthorID: 850,
					Ctime:    123,
					Utime:    234,
				}
				err := s.db.Create(&article).Error
				assert.NoError(t, err)
			},

			after: func(t *testing.T) {
				// 由于是别人的帖子，修改应该是失败的，因此id=15的帖子的数据是不变的。

				var article dao.Article
				var authorID int64 = 850
				var articleID int64 = 15

				err := s.db.Where("id", articleID).First(&article).Error
				assert.NoError(t, err)

				assert.Equal(t, dao.Article{
					ID:       articleID,
					Title:    "my article",
					Content:  "my content",
					Status:   domain.ArticleStatusPublished,
					Ctime:    123,
					Utime:    234,
					AuthorID: authorID,
				}, article)
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

func (s *ArticleTestSuite) TestPublish() {
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
			name: "新建帖子，发表帖子",
			article: Article{
				Title:   "new title",
				Content: "new content",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "ok",
				Data: 1,
			},
			reqBuilder: func(t *testing.T, article Article) *http.Request {
				buf, err := json.Marshal(article)
				assert.NoError(t, err)

				// 构建起请求
				req, err := http.NewRequest("POST", "/articles/publish", bytes.NewBuffer(buf))
				assert.NoError(t, err)

				req.Header.Set("Content-Type", "application/json; charset=utf-8")
				return req
			},
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				var article dao.Article
				err := s.db.First(&article).Error
				assert.NoError(t, err)
				assert.True(t, article.ID > 0)
				assert.True(t, article.Ctime > 0)
				assert.True(t, article.Utime > 0)
				article.ID = 0
				article.Ctime = 0
				article.Utime = 0

				expectedArt := dao.Article{
					Title:    "new title",
					Content:  "new content",
					AuthorID: 1000,
					Status:   domain.ArticleStatusPublished,
				}
				assert.Equal(t, expectedArt, article)

				var pubArt dao.PublishArticle
				err = s.db.First(&pubArt).Error
				assert.NoError(t, err)
				assert.True(t, pubArt.ID > 0)
				assert.True(t, pubArt.Ctime > 0)
				assert.True(t, pubArt.Utime > 0)
				pubArt.ID = 0
				pubArt.Ctime = 0
				pubArt.Utime = 0

				assert.Equal(t, dao.PublishArticle(expectedArt), pubArt)
			},
		},
		{
			name: "编辑帖子，首次发表",
			article: Article{
				ID:      15,
				Title:   "new title",
				Content: "new content",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "ok",
				Data: 15,
			},
			reqBuilder: func(t *testing.T, article Article) *http.Request {
				buf, err := json.Marshal(article)
				assert.NoError(t, err)

				// 构建起请求
				req, err := http.NewRequest("POST", "/articles/publish", bytes.NewBuffer(buf))
				assert.NoError(t, err)

				req.Header.Set("Content-Type", "application/json; charset=utf-8")
				return req
			},
			before: func(t *testing.T) {
				// 制作库插入一条数据
				err := s.db.Create(&dao.Article{
					ID:       15,
					Title:    "my article",
					Content:  "my content",
					Status:   domain.ArticleStatusUnpublished,
					AuthorID: 1000,
					Ctime:    123,
					Utime:    456,
				}).Error
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				var article dao.Article
				err := s.db.Where("id=?", 15).First(&article).Error
				assert.NoError(t, err)

				assert.True(t, article.Utime > 456)
				article.Utime = 0

				expectedArt := dao.Article{
					ID:       15,
					Title:    "new title",
					Content:  "new content",
					AuthorID: 1000,
					Status:   domain.ArticleStatusPublished,
					Ctime:    123,
				}
				assert.Equal(t, expectedArt, article)

				var pubArt dao.PublishArticle
				err = s.db.Where("id = ?", 15).First(&pubArt).Error
				assert.NoError(t, err)
				assert.True(t, pubArt.Ctime > 0)
				assert.True(t, pubArt.Utime > 0)
				pubArt.Ctime = 0
				pubArt.Utime = 0
				expectedArt.Ctime = 0

				assert.Equal(t, dao.PublishArticle(expectedArt), pubArt)
			},
		},
		{
			name: "编辑帖子，再次发表",
			article: Article{
				ID:      20,
				Title:   "new title",
				Content: "new content",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "ok",
				Data: 20,
			},
			reqBuilder: func(t *testing.T, article Article) *http.Request {
				buf, err := json.Marshal(article)
				assert.NoError(t, err)

				// 构建起请求
				req, err := http.NewRequest("POST", "/articles/publish", bytes.NewBuffer(buf))
				assert.NoError(t, err)

				req.Header.Set("Content-Type", "application/json; charset=utf-8")
				return req
			},
			before: func(t *testing.T) {
				// 制作库插入一条数据
				err := s.db.Create(&dao.Article{
					ID:       20,
					Title:    "my article",
					Content:  "my content",
					Status:   domain.ArticleStatusUnpublished,
					AuthorID: 1000,
					Ctime:    123,
					Utime:    456,
				}).Error
				require.NoError(t, err)
				// 线上库插入一条数据
				err = s.db.Create(&dao.PublishArticle{
					ID:       20,
					Title:    "my article",
					Content:  "my content",
					Status:   domain.ArticleStatusPublished,
					AuthorID: 1000,
					Ctime:    123,
					Utime:    456,
				}).Error
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				var article dao.Article
				err := s.db.Where("id=?", 20).First(&article).Error
				assert.NoError(t, err)

				assert.True(t, article.Utime > 456)
				article.Utime = 0

				expectedArt := dao.Article{
					ID:       20,
					Title:    "new title",
					Content:  "new content",
					AuthorID: 1000,
					Status:   domain.ArticleStatusPublished,
					Ctime:    123,
				}
				assert.Equal(t, expectedArt, article)

				var pubArt dao.PublishArticle
				err = s.db.Where("id = ?", 20).First(&pubArt).Error
				assert.NoError(t, err)
				assert.True(t, pubArt.Utime > 0)
				pubArt.Utime = 0

				assert.Equal(t, dao.PublishArticle(expectedArt), pubArt)
			},
		},
		{
			name: "发表别人帖子",
			article: Article{
				ID:      22,
				Title:   "change title",
				Content: "change content",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "failed",
				Code: 5,
			},
			reqBuilder: func(t *testing.T, article Article) *http.Request {
				buf, err := json.Marshal(article)
				assert.NoError(t, err)

				// 构建起请求
				req, err := http.NewRequest("POST", "/articles/publish", bytes.NewBuffer(buf))
				assert.NoError(t, err)

				req.Header.Set("Content-Type", "application/json; charset=utf-8")
				return req
			},
			before: func(t *testing.T) {
				// 制作库插入一条数据
				err := s.db.Create(&dao.Article{
					ID:       22,
					Title:    "hi",
					Content:  "welcome back",
					Status:   domain.ArticleStatusPublished,
					AuthorID: 900,
					Ctime:    123,
					Utime:    456,
				}).Error
				require.NoError(t, err)
				// 线上库插入一条数据
				err = s.db.Create(&dao.PublishArticle{
					ID:       22,
					Title:    "hi",
					Content:  "welcome back",
					Status:   domain.ArticleStatusPublished,
					AuthorID: 900,
					Ctime:    123,
					Utime:    456,
				}).Error
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				var article dao.Article
				err := s.db.Where("id=?", 22).First(&article).Error
				assert.NoError(t, err)

				expectedArt := dao.Article{
					ID:       22,
					Title:    "hi",
					Content:  "welcome back",
					AuthorID: 900,
					Status:   domain.ArticleStatusPublished,
					Ctime:    123,
					Utime:    456,
				}
				assert.Equal(t, expectedArt, article)

				var pubArt dao.PublishArticle
				err = s.db.Where("id = ?", 22).First(&pubArt).Error
				assert.NoError(t, err)

				assert.Equal(t, dao.PublishArticle(expectedArt), pubArt)
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

func (s *ArticleTestSuite) TestArticleList() {
	t := s.T()

	testCases := []struct {
		name string

		req        web.ListReq
		reqBuilder func(t *testing.T, req web.ListReq) *http.Request
		before     func(t *testing.T)
		after      func(t *testing.T)

		// 期望的输出
		wantCode int
		wantRes  ginx.Result
	}{
		{
			name: "查询接口测试",
			req: web.ListReq{
				Offset: 0,
				Limit:  10,
			},
			wantCode: http.StatusOK,
			wantRes: ginx.Result{
				Msg: "ok",
				Data: []web.ArticleVO{
					{
						ID:      1,
						Title:   "title1",
						Content: "content1",
						CTime:   "2024-10-01 11:00:10",
						UTime:   "2024-10-02 12:00:20",
					},
					{
						ID:      2,
						Title:   "title2",
						Content: "content2",
						CTime:   "2024-05-01 11:00:10",
						UTime:   "2024-05-02 12:00:20",
					},
				},
			},
			reqBuilder: func(t *testing.T, listReq web.ListReq) *http.Request {
				//buf, err := json.Marshal(listReq)
				//assert.NoError(t, err)

				params := url.Values{}
				params.Add("limit", strconv.Itoa(listReq.Limit))
				params.Add("offset", strconv.Itoa(listReq.Offset))
				// 构建起请求
				req, err := http.NewRequest("GET", "/articles/list?"+params.Encode(), nil)
				assert.NoError(t, err)
				// req.Header.Set("Content-Type", "application/json; charset=utf-8")
				return req
			},
			before: func(t *testing.T) {
				// 插入一些测试数据
				c1, err := time.Parse(time.DateTime, "2024-10-01 11:00:10")
				assert.NoError(t, err)
				u1, err := time.Parse(time.DateTime, "2024-10-02 12:00:20")
				assert.NoError(t, err)

				c2, err := time.Parse(time.DateTime, "2024-05-01 11:00:10")
				assert.NoError(t, err)
				u2, err := time.Parse(time.DateTime, "2024-05-02 12:00:20")
				assert.NoError(t, err)

				arts := []dao.Article{
					{
						ID:       3,
						Title:    "title1",
						Content:  "content1",
						Ctime:    c1.UnixMilli(),
						Utime:    u1.UnixMilli(),
						AuthorID: 1000,
					},
					{
						ID:       4,
						Title:    "title2",
						Content:  "content2",
						Ctime:    c2.UnixMilli(),
						Utime:    u2.UnixMilli(),
						AuthorID: 1000,
					},
				}
				s.db.Create(&arts)
			},
			after: func(t *testing.T) {
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.before(t)

			// 1. 构建请求
			req := testCase.reqBuilder(t, testCase.req)

			// 2. 处理该请求并写入响应
			resp := httptest.NewRecorder()
			s.server.ServeHTTP(resp, req)

			// 3. 校验响应
			assert.Equal(t, testCase.wantCode, resp.Code)

			body := resp.Body.Bytes()
			var webResult ginx.Result
			err := json.Unmarshal(body, &webResult)
			require.NoError(t, err)

			assert.Equal(t, testCase.wantRes.Msg, webResult.Msg)

			arts, ok := webResult.Data.([]interface{})
			require.True(t, ok)
			assert.True(t, len(arts) > 0)

			testCase.after(t)
		})
	}
}

// 测试方法的运行入口
func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
}
