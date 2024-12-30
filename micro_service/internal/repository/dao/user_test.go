package dao

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGORMUserDao_Insert(t *testing.T) {
	testCases := []struct {
		name string

		mock func(t *testing.T) *sql.DB

		ctx  context.Context
		user User

		wantErr error
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				// mock出sql.DB的实现
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				// 预期要执行的行为
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnError(nil)
				return db
			},
			ctx: context.Background(),
			user: User{
				Email: sql.NullString{
					String: "test@qq.com",
					Valid:  true,
				},
				Birthday: 0,
			},
			wantErr: nil,
		},
		{
			name: "重复账号",
			mock: func(t *testing.T) *sql.DB {
				// mock出sql.DB的实现
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				// 预期要执行的行为
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnError(&mysqlDriver.MySQLError{Number: 1062})
				return db
			},
			ctx: context.Background(),
			user: User{
				Email: sql.NullString{
					String: "test@qq.com",
					Valid:  true,
				},
				Birthday: 0,
			},
			wantErr: ErrUserDuplicate,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.mock(t)

			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})

			require.NoError(t, err)

			ud := NewUserDao(db)
			err = ud.Insert(tc.ctx, tc.user)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
