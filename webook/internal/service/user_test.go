package service

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/domain"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository"
	repomocks "Neo/Workplace/goland/src/GeekGo/webook/internal/repository/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func TestUserServiceImpl_Login(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository
		//输入
		//ctx  context.Context
		user domain.User

		//输出
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "login success",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "test@gmail.com").
					Return(domain.User{
						Email:    "test@gmail.com",
						Password: "$2a$10$f9LW51lHG8OM2Zk8xFuAZO9aaycuNDEdEIm82o0BbumTkdaMMVPSm",
						Phone:    "13837111118",
						Ctime:    now,
					}, nil)
				return repo
			},
			user: domain.User{
				Email:    "test@gmail.com",
				Password: "12345678",
			},

			wantUser: domain.User{
				Email:    "test@gmail.com",
				Password: "$2a$10$f9LW51lHG8OM2Zk8xFuAZO9aaycuNDEdEIm82o0BbumTkdaMMVPSm",
				Phone:    "13837111118",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "test@gmail.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			user: domain.User{
				Email:    "test@gmail.com",
				Password: "12345678",
			},

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "DB错误（系统错误）",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "test@gmail.com").
					Return(domain.User{}, errors.New("mock DB error"))
				return repo
			},
			user: domain.User{
				Email:    "test@gmail.com",
				Password: "12345678",
			},

			wantUser: domain.User{},
			wantErr:  errors.New("mock DB error"),
		},
		{
			name: "密码不对",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "test@gmail.com").
					Return(domain.User{
						Email:    "test@gmail.com",
						Password: "$2a$10$f9LW51lHG8OM2Zk8xFuAZO9aaycuNDEdEIm82o0BbumTkdaMMVPSm",
						Phone:    "13837111118",
						Ctime:    now,
					}, nil)
				return repo
			},
			user: domain.User{
				Email:    "test@gmail.com",
				Password: "23456781",
			},

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewUserService(tc.mock(ctrl))
			u, err := svc.Login(context.Background(), tc.user)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}

func TestEncrypted(t *testing.T) {
	res, err := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
	if err == nil {
		t.Log(string(res))
	}
}
