package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/repository"
	"learn_go/webook/pkg/logger"
)

var (
	ErrDuplicateUser          = repository.ErrDuplicateUser
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")

	ErrUserNotFound = repository.ErrUserNotFound
)

type UserService interface {
	Login(ctx context.Context, email, password string) (domain.User, error)
	SignUp(ctx context.Context, u domain.User) error

	UpdateNonSensitiveInfo(ctx context.Context, u domain.User) error
	Profile(ctx context.Context, id int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain.User, error)
}

type userService struct {
	userRepo repository.UserRepository
	logger   logger.LoggerV2
}

func NewUserService(userRepository repository.UserRepository, logger logger.LoggerV2) UserService {
	return &userService{
		userRepo: userRepository,
		logger:   logger,
	}
}

func (svc *userService) Login(ctx context.Context, email, password string) (domain.User, error) {
	user, err := svc.userRepo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidEmailOrPassword
	}
	return user, nil
}

func (svc *userService) SignUp(ctx context.Context, u domain.User) error {
	// 对密码进行加密
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.userRepo.Create(ctx, u)
}

func (svc *userService) UpdateNonSensitiveInfo(ctx context.Context, u domain.User) error {
	// note：我的惯例做法是先找出该对象，然后再进行更新。
	// 其实究竟是否需要确认该对象是否存在，看具体需求吧。比如对更新的资源的权限控制?
	return svc.userRepo.UpdateNonZeroFields(ctx, u)
}

func (svc *userService) Profile(ctx context.Context, id int64) (domain.User, error) {
	return svc.userRepo.FindByID(ctx, id)
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	user, err := svc.userRepo.FindByPhone(ctx, phone)

	// err=nil或者err=其他错误
	if err != ErrUserNotFound {
		return user, err
	}

	// 未找到用户的处理
	err = svc.userRepo.Create(ctx, domain.User{
		Phone: phone,
	})
	if err != nil && err != repository.ErrDuplicateUser {
		return domain.User{}, err
	}

	svc.logger.Info("手机号重复了")

	// 手机号存在重复的处理
	return svc.userRepo.FindByPhone(ctx, phone)

	//// 未找到用户的处理
	//if err == ErrUserNotFound {
	//	newUser := domain.User{
	//		Phone: phone,
	//	}
	//	// 问题1：返回对象的时间是否需要更新?
	//	err = svc.userRepo.Create(ctx, newUser)
	//	if err != nil {
	//		return domain.User{}, err
	//	}
	//	return newUser, nil
	//}
	//
	//return user, err
}

func (svc *userService) FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain.User, error) {
	user, err := svc.userRepo.FindByOpenId(ctx, wechatInfo.OpenId)

	// err=nil或者err=其他错误
	if err != ErrUserNotFound {
		return user, err
	}

	// 未找到用户的处理
	err = svc.userRepo.Create(ctx, domain.User{
		WechatInfo: wechatInfo,
	})
	if err != nil && err != repository.ErrDuplicateUser {
		return domain.User{}, err
	}

	// 手机号存在重复的处理
	return svc.userRepo.FindByOpenId(ctx, wechatInfo.OpenId)
}
