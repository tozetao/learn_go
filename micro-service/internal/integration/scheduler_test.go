package integration

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/integration/startup"
	"learn_go/webook/internal/job"
	"learn_go/webook/internal/repository/dao"
	"net/http"
	"testing"
	"time"
)

type SchedulerTestSuite struct {
	suite.Suite
	db *gorm.DB
}

// hook, 在测试启动之前触发
func (s *SchedulerTestSuite) SetupSuite() {
	s.db = startup.NewDB()

	now := time.Now()
	testJob := dao.Job{
		Name:     "test_job",
		Executor: "executor:local",

		Expression: "*/10 * * * * ?",
		Status:     1,
		Cfg:        "this is a test job",
		CTime:      now.UnixMilli(),
		UTime:      now.UnixMilli(),
	}

	err := s.db.Create(&testJob).Error
	require.NoError(s.T(), err)
}

func (s *SchedulerTestSuite) TearDownSuite() {
	//t := s.T()
	//err := s.db.Exec("truncate table `jobs`").Error
	//assert.NoError(t, err)
}

func (s *SchedulerTestSuite) TestJob() {
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
		{},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.before(t)

		})
	}

}

func TestJob1(t *testing.T) {
	db := startup.NewDB()
	err := db.Exec("truncate table `jobs`").Error
	require.NoError(t, err)

	now := time.Now()
	testJob := dao.Job{
		Name:     "test_job",
		Executor: "executor:local",

		Expression: "*/5 * * * * ?",
		Status:     1,
		Cfg:        "this is a test job",
		CTime:      now.UnixMilli(),
		UTime:      now.UnixMilli(),
	}

	err = db.Create(&testJob).Error
	require.NoError(t, err)
}

func TestJob2(t *testing.T) {
	executor := job.NewLocalExecutor()
	executor.RegisterFunc("test_job", func(ctx context.Context, j domain.Job) error {
		t.Logf("job id: %d, name: %s", j.ID, j.Name)
		return nil
	})

	scheduler := startup.InitScheduler()

	// 注册executor
	scheduler.Register(executor.Name(), executor)

	// 开始执行
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	scheduler.Schedule(ctx)
}
