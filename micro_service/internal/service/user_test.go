package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/repository"
	repomocks "learn_go/webook/internal/repository/mocks"
	"learn_go/webook/pkg/logger"
	"testing"
)

func Test_userService_Login(t *testing.T) {
	testCases := []struct {
		name string

		// mock依赖对象
		userRepo func(ctrl *gomock.Controller) repository.UserRepository

		ctx      context.Context
		email    string
		password string

		wantErr  error
		wantUser domain.User
	}{
		{
			name: "登录成功",

			userRepo: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(context.Background(), "test@qq.com").
					Return(domain.User{
						ID:       1,
						Email:    "test@qq.com",
						Password: "$2a$10$iJ98au4JWA0kpxGHKXbyxOrNO2XABQ5G7tNX1lbnSOSY695QjCTri",
					}, nil)

				return repo
			},

			ctx:      context.Background(),
			email:    "test@qq.com",
			password: "123456",

			wantErr: nil,
			wantUser: domain.User{
				ID:       1,
				Email:    "test@qq.com",
				Password: "$2a$10$iJ98au4JWA0kpxGHKXbyxOrNO2XABQ5G7tNX1lbnSOSY695QjCTri",
			},
		},
		{
			name: "用户不存在",

			userRepo: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(context.Background(), "test@qq.com").
					Return(domain.User{}, repository.ErrUserNotFound)

				return repo
			},

			ctx:      context.Background(),
			email:    "test@qq.com",
			password: "123456",

			wantErr:  ErrUserNotFound,
			wantUser: domain.User{},
		},
		{
			name: "其他错误，比如DB错误",

			userRepo: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(context.Background(), "test@qq.com").
					Return(domain.User{}, errors.New("db超时"))

				return repo
			},

			ctx:      context.Background(),
			email:    "test@qq.com",
			password: "123456",

			wantErr:  errors.New("db超时"),
			wantUser: domain.User{},
		},
		{
			name: "密码错误",
			userRepo: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(context.Background(), "test@qq.com").
					Return(domain.User{}, nil)

				return repo
			},
			ctx:      context.Background(),
			email:    "test@qq.com",
			password: "qq123456",
			wantErr:  ErrInvalidEmailOrPassword,
			wantUser: domain.User{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.userRepo(ctrl)

			svc := NewUserService(repo, logger.NewNopLogger())

			gotUser, err := svc.Login(tc.ctx, tc.email, tc.password)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, gotUser)
		})
	}
}

func TestPassword(t *testing.T) {
	password := "123456"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err == nil {
		t.Log(string(hash))
	}
}
