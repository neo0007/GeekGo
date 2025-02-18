package service

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/domain"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicate = repository.ErrUserDuplicate
var ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")

type UserService interface {
	Signup(c context.Context, u domain.User) error
	Login(c context.Context, u domain.User) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain.User, error)
}

type UserServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{repo: repo}
}

func (svc *UserServiceImpl) Signup(c context.Context, u domain.User) error {
	//考虑加密
	//考虑存储
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)

	return svc.repo.Create(c, u)
}

func (svc *UserServiceImpl) Login(c context.Context, u domain.User) (domain.User, error) {
	//先找用户
	ud, err := svc.repo.FindByEmail(c, u.Email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	//比较密码
	if err := bcrypt.CompareHashAndPassword([]byte(ud.Password), []byte(u.Password)); err != nil {
		//记录 debug 日志
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return ud, nil
}

func (svc *UserServiceImpl) Profile(ctx context.Context, id int64) (domain.User, error) {
	u, err := svc.repo.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (svc *UserServiceImpl) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err == nil {
		return u, nil
	}
	if err != repository.ErrUserNotFound {
		// err 为 nil 会进来这里
		// 不为 ErrUserNotFound 也会进来这里
		return domain.User{}, err
	}
	u = domain.User{
		Phone: phone,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil && !errors.Is(err, ErrUserDuplicate) {
		return u, err
	}
	// 这里会遇到主从延迟的问题，如果按照下面代码
	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *UserServiceImpl) FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error) {
	u, err := svc.repo.FindByWechat(ctx, info.OpenID)
	if err == nil {
		return u, nil
	}
	if err != repository.ErrUserNotFound {
		// err 不为 ErrUserNotFound 会进来这里
		return domain.User{}, err
	}
	u = domain.User{
		WechatInfo: info,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil && !errors.Is(err, ErrUserDuplicate) {
		return u, err
	}
	// 这里会遇到主从延迟的问题，如果按照下面代码
	return svc.repo.FindByWechat(ctx, info.OpenID)
}
