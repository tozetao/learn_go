package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicate  = errors.New("duplicate account")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

//go:generate mockgen -source=./user.go -package=daomocks -destination=./mocks/user.mock.go UserDao
type UserDao interface {
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindByID(ctx context.Context, id int64) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByOpenId(ctx context.Context, openId string) (User, error)
	Insert(ctx context.Context, u User) error
	UpdateById(ctx context.Context, u User) error
}

type GORMUserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDao {
	return &GORMUserDao{
		db: db,
	}
}

func (dao *GORMUserDao) FindByID(ctx context.Context, id int64) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("`id` = ?", id).First(&user).Error
	return user, err
}

func (dao *GORMUserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}

func (dao *GORMUserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	return user, err
}

func (dao *GORMUserDao) FindByOpenId(ctx context.Context, openId string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("open_id = ?", openId).First(&user).Error
	return user, err
}

func (dao *GORMUserDao) Insert(ctx context.Context, u User) error {
	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			return ErrUserDuplicate
		}
	}
	return err
}

func (dao *GORMUserDao) UpdateById(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	m := map[string]any{
		"u_time":   now,
		"nickname": u.Nickname,
		"birthday": u.Birthday,
		"about_me": u.AboutMe,
	}
	// Model: 绑定要更新的对象。如果对象存在主键ID，则自动添加主键ID的where条件。
	// Updates的参数的结构体时，默认会忽略nil值字段的更新。
	return dao.db.WithContext(ctx).Model(&u).Where("id=?", u.ID).Updates(m).Error
}

// User PO(persistent object), entity, model
type User struct {
	ID       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string

	Nickname string `gorm:"type=varchar(30)"`
	// YYYY-MM-DD
	Birthday int64
	AboutMe  string `gorm:"type=varchar(255)"`

	// 表示字段可以为null
	Phone sql.NullString `gorm:"unique"`

	UnionId sql.NullString `gorm:"unique"`
	OpenId  sql.NullString `gorm:"unique"`

	CTime int64
	UTime int64
}
