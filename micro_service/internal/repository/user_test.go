package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/repository/cache"
	cachemocks "learn_go/webook/internal/repository/cache/mocks"
	"learn_go/webook/internal/repository/dao"
	daomocks "learn_go/webook/internal/repository/dao/mocks"
	"testing"
	"time"
)

func Test_userRepository_FindByID(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name string

		ctx context.Context
		id  int64

		mock func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)

		wantErr  error
		wantUser domain.User
	}{
		{
			name: "缓存未命中, 查询成功.",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				userDao := daomocks.NewMockUserDao(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)

				id := int64(1)

				userCache.EXPECT().Get(gomock.Any(), id).Return(domain.User{}, ErrKeyNotExist)

				userDao.EXPECT().FindByID(gomock.Any(), id).Return(dao.User{
					ID: 1,
					Email: sql.NullString{
						String: "test@test.com",
						Valid:  true,
					},
					Password: "123456",
					CTime:    now.UnixMilli(),
				}, nil)

				userCache.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				return userDao, userCache
			},
			ctx: context.Background(),
			id:  1,

			wantErr: nil,
			wantUser: domain.User{
				ID:       1,
				Email:    "test@test.com",
				Password: "123456",
				Birthday: time.UnixMilli(0),
				Ctime:    time.UnixMilli(now.UnixMilli()),
			},
		},
		{
			name: "缓存命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				userDao := daomocks.NewMockUserDao(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)

				id := int64(1)

				userCache.EXPECT().Get(gomock.Any(), id).Return(domain.User{
					ID:       1,
					Email:    "test@test.com",
					Password: "123456",
					Birthday: time.UnixMilli(0),
					Ctime:    time.UnixMilli(now.UnixMilli()),
				}, nil)

				return userDao, userCache
			},
			ctx: context.Background(),
			id:  1,

			wantErr: nil,
			wantUser: domain.User{
				ID:       1,
				Email:    "test@test.com",
				Password: "123456",
				Birthday: time.UnixMilli(0),
				Ctime:    time.UnixMilli(now.UnixMilli()),
			},
		},
		{
			name: "缓存异常",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				userDao := daomocks.NewMockUserDao(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)

				id := int64(1)

				userCache.EXPECT().Get(gomock.Any(), id).Return(domain.User{}, errors.New("cache error"))
				return userDao, userCache
			},
			ctx: context.Background(),
			id:  1,

			wantErr:  errors.New("cache error"),
			wantUser: domain.User{},
		},
		{
			name: "找不到用户",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				userDao := daomocks.NewMockUserDao(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)

				id := int64(1)

				userCache.EXPECT().Get(gomock.Any(), id).Return(domain.User{}, ErrKeyNotExist)

				userDao.EXPECT().FindByID(gomock.Any(), id).Return(dao.User{}, dao.ErrRecordNotFound)

				return userDao, userCache
			},
			ctx: context.Background(),
			id:  1,

			wantErr:  dao.ErrRecordNotFound,
			wantUser: domain.User{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userDao, userCache := tc.mock(ctrl)

			repo := NewUserRepository(userDao, userCache)

			user, err := repo.FindByID(tc.ctx, tc.id)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}

}

func TestTime(t *testing.T) {
	var birthday int64 = 0
	b := time.UnixMilli(birthday)
	fmt.Printf("%v\n%v\n", b, b.UnixMilli())

	inner := struct {
		birthday time.Time
	}{}
	fmt.Printf("%v\n%v\n%v\n", inner.birthday, inner.birthday.UnixMilli(), inner.birthday.IsZero())
}
