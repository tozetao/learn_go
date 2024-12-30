package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/integration/startup"
	"learn_go/webook/internal/repository/dao"
	"learn_go/webook/internal/web"
	"net/http"
	"net/http/httptest"
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
type ArticleTestSuiteV1 struct {
	suite.Suite
	server          *gin.Engine
	db              *mongo.Database
	artCol          *mongo.Collection
	publishedArtCol *mongo.Collection
}

// 测试方法的运行入口
func TestMangoArticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuiteV1{})
}

// hook, 在测试启动之前触发
func (s *ArticleTestSuiteV1) SetupSuite() {
	// 初始化mangodb数据库
	s.db = startup.NewMangoDB()
	s.artCol = s.db.Collection("articles")
	s.publishedArtCol = s.db.Collection("published_articles")
	// 清空原有表数据
	s.truncate()

	s.server = gin.Default()
	s.server.Use(func(ctx *gin.Context) {
		ctx.Set("user", &web.UserClaims{
			Uid: 1000,
		})
	})
	// 新建handler
	node, err := snowflake.NewNode(1)
	assert.NoError(s.T(), err)
	handler := startup.InitArticleHandlerV1(dao.NewMongoArticleDao(s.db, node))
	// 注册路由
	handler.RegisterRoutes(s.server)
}

func (s *ArticleTestSuiteV1) TearDownSuite() {

}

func (s *ArticleTestSuiteV1) truncate() {
	t := s.T()

	delRes, err := s.artCol.DeleteMany(context.Background(), bson.D{})
	assert.NoError(t, err)
	t.Log("deleted articles: ", delRes.DeletedCount)

	pubArtDelRes, err := s.publishedArtCol.DeleteMany(context.Background(), bson.D{})
	assert.NoError(t, err)
	t.Log("deleted published articles: ", pubArtDelRes.DeletedCount)
}

func (s *ArticleTestSuiteV1) TestEdit() {
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
				Msg: "ok",
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
				// 制作库的比较
				var article dao.Article
				filter := bson.M{"title": "hello"}
				err := s.artCol.FindOne(context.Background(), filter).Decode(&article)
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
				_, err := s.artCol.InsertOne(context.Background(), article)
				assert.NoError(t, err)
			},

			after: func(t *testing.T) {
				var article dao.Article
				err := s.artCol.FindOne(context.Background(), bson.M{"id": 5}).Decode(&article)
				assert.NoError(t, err)

				// 更新后的时间必定会大于原本的UTime
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
				_, err := s.artCol.InsertOne(context.TODO(), article)
				assert.NoError(t, err)
			},

			after: func(t *testing.T) {
				// 由于是别人的帖子，修改应该是失败的，因此id=15的帖子的数据是不变的。
				var authorID int64 = 850
				var articleID int64 = 15
				var article dao.Article

				err := s.artCol.FindOne(context.Background(), bson.M{"id": articleID}).Decode(&article)
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
			// assert.Equal(t, testCase.wantRes.Data, webResult.Data)

			testCase.after(t)
		})
	}
}

func (s *ArticleTestSuiteV1) TestPublish() {
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
				// 制作库的比较
				var article dao.Article
				filter := bson.M{"title": "new title"}
				err := s.artCol.FindOne(context.Background(), filter).Decode(&article)
				assert.NoError(t, err)
				articleID := article.ID

				assert.True(t, article.ID > 0)
				assert.True(t, article.Ctime > 0)
				assert.True(t, article.Utime > 0)
				article.ID = 0
				article.Ctime = 0
				article.Utime = 0
				assert.Equal(t, dao.Article{
					Title:    "new title",
					Content:  "new content",
					AuthorID: 1000,
					Status:   domain.ArticleStatusPublished,
				}, article)

				// 制作库的比较
				// bson.M{"article.title": "hello"}
				var pubArt dao.PublishArticle
				err = s.publishedArtCol.FindOne(context.Background(), bson.D{
					{"title", "new title"},
				}).Decode(&pubArt)
				assert.NoError(t, err)
				assert.Equal(t, pubArt.ID, articleID)
				assert.True(t, pubArt.Ctime > 0)
				assert.True(t, pubArt.Utime > 0)
				pubArt.ID = 0
				pubArt.Ctime = 0
				pubArt.Utime = 0
				assert.Equal(t, dao.PublishArticle{
					Title:    "new title",
					Content:  "new content",
					AuthorID: 1000,
					Status:   domain.ArticleStatusPublished,
				}, pubArt)
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
			before: func(t *testing.T) {
				// 制作库插入一条数据
				art := dao.Article{
					ID:       15,
					Title:    "my article",
					Content:  "my content",
					Status:   domain.ArticleStatusUnpublished,
					AuthorID: 1000,
					Ctime:    123,
					Utime:    456,
				}
				_, err := s.artCol.InsertOne(context.Background(), art)
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				var article dao.Article
				err := s.artCol.FindOne(context.Background(), bson.M{"id": 15}).Decode(&article)
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
				err = s.publishedArtCol.FindOne(context.Background(), bson.M{"id": 15}).Decode(&pubArt)
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
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				// 制作库插入一条数据
				art := dao.Article{
					ID:       20,
					Title:    "my article",
					Content:  "my content",
					Status:   domain.ArticleStatusUnpublished,
					AuthorID: 1000,
					Ctime:    123,
					Utime:    456,
				}
				_, err := s.artCol.InsertOne(ctx, art)
				require.NoError(t, err)
				// 线上库插入一条数据
				_, err = s.publishedArtCol.InsertOne(ctx, dao.PublishArticle(art))
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				var article dao.Article
				err := s.artCol.FindOne(context.Background(), bson.M{"id": 20}).Decode(&article)
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
				err = s.publishedArtCol.FindOne(context.Background(), bson.M{"id": 20}).Decode(&pubArt)
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
				art := dao.Article{
					ID:       22,
					Title:    "hi",
					Content:  "welcome back",
					Status:   domain.ArticleStatusPublished,
					AuthorID: 900,
					Ctime:    123,
					Utime:    456,
				}
				_, err := s.artCol.InsertOne(context.Background(), art)
				require.NoError(t, err)

				// 线上库插入一条数据
				_, err = s.publishedArtCol.InsertOne(context.Background(), dao.PublishArticle(art))
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				var article dao.Article
				err := s.artCol.FindOne(context.Background(), bson.M{"id": 22}).Decode(&article)
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
				err = s.publishedArtCol.FindOne(context.Background(), bson.M{"id": 22}).Decode(&pubArt)
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
			if testCase.wantRes.Data > 0 {
				assert.True(t, webResult.Data > 0)
			}

			testCase.after(t)
		})
	}
}
