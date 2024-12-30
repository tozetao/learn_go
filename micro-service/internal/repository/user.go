package repository

import (
	"context"
	"database/sql"
	"github.com/redis/go-redis/v9"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/repository/cache"
	"learn_go/webook/internal/repository/dao"
	"log"
	"time"
)

var (
	ErrDuplicateUser = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrRecordNotFound
	ErrKeyNotExist   = redis.Nil
)

type UserRepository interface {
	FindByID(ctx context.Context, id int64) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByOpenId(ctx context.Context, openId string) (domain.User, error)
	UpdateNonZeroFields(ctx context.Context, u domain.User) error
	Create(ctx context.Context, u domain.User) error
}

type userRepository struct {
	ud    dao.UserDao
	cache cache.UserCache
}

func NewUserRepository(ud dao.UserDao, c cache.UserCache) UserRepository {
	return &userRepository{
		ud:    ud,
		cache: c,
	}
}

func (repo *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.ud.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return toDomain(u), err
}

func (repo *userRepository) Create(ctx context.Context, u domain.User) error {
	now := time.Now().UnixMilli()

	entity := toEntity(u)
	entity.CTime = now
	entity.UTime = now

	return repo.ud.Insert(ctx, entity)
}

// FindByID error的几种情况：error != nil时表示找到用户；error == ErrKeyNotExist表示找不到用户，其他情况都是错误。
func (repo *userRepository) FindByID(ctx context.Context, id int64) (domain.User, error) {
	user, err := repo.cache.Get(ctx, id)

	// 缓存存在数据
	if err == nil {
		return user, nil
	}

	// key不存在会从数据库进行查找
	if err == ErrKeyNotExist {
		userModel, err := repo.ud.FindByID(ctx, id)
		if err != nil {
			return domain.User{}, err
		}
		err = repo.cache.Set(ctx, toDomain(userModel), time.Minute*30)
		if err != nil {
			// 打印日志并告警
			log.Println("cache.Set err: ", err)
		}
		return toDomain(userModel), nil
	}

	// 缓存异常
	// 对于这种情况，如果仍然要继续去数据库进行查询，要限流，保证数据库可用性。同时记录错误日志。这里忽略，简单处理
	return domain.User{}, err
}

func (repo *userRepository) UpdateNonZeroFields(ctx context.Context, u domain.User) error {
	// note: 个人感觉应该过滤nil字段，才符合该方法名的定义。
	return repo.ud.UpdateById(ctx, toEntity(u))
}

func (repo *userRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	user, err := repo.ud.FindByPhone(ctx, phone)
	return toDomain(user), err
}

func (repo *userRepository) FindByOpenId(ctx context.Context, openId string) (domain.User, error) {
	user, err := repo.ud.FindByOpenId(ctx, openId)
	return toDomain(user), err
}

func toDomain(user dao.User) domain.User {
	return domain.User{
		ID:       user.ID,
		Email:    user.Email.String,
		Phone:    user.Phone.String,
		Password: user.Password,
		Nickname: user.Nickname,
		AboutMe:  user.AboutMe,
		Birthday: time.UnixMilli(user.Birthday),
		Ctime:    time.UnixMilli(user.CTime),
		WechatInfo: domain.WechatInfo{
			UnionId: user.UnionId.String,
			OpenId:  user.OpenId.String,
		},
	}
}

func toEntity(d domain.User) dao.User {
	return dao.User{
		ID: d.ID,
		Email: sql.NullString{
			String: d.Email,
			Valid:  d.Email != "",
		},
		Phone: sql.NullString{
			String: d.Phone,
			Valid:  d.Phone != "",
		},
		UnionId: sql.NullString{
			String: d.WechatInfo.UnionId,
			Valid:  d.WechatInfo.UnionId != "",
		},
		OpenId: sql.NullString{
			String: d.WechatInfo.OpenId,
			Valid:  d.WechatInfo.OpenId != "",
		},
		Password: d.Password,
		Nickname: d.Nickname,
		Birthday: d.Birthday.UnixMilli(),
		AboutMe:  d.AboutMe,
	}
}
