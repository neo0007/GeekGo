package service

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/domain"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
var ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) Signup(c context.Context, u domain.User) error {
	//考虑加密
	//考虑存储
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)

	return svc.repo.Create(c, u)
}

func (svc *UserService) Login(c context.Context, u domain.User) (domain.User, error) {
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
