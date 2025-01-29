package repository

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/domain"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository/cache"
	cachemocks "Neo/Workplace/goland/src/GeekGo/webook/internal/repository/cache/mocks"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository/cache/redis"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository/dao"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository/dao/entity"
	daomocks "Neo/Workplace/goland/src/GeekGo/webook/internal/repository/dao/mocks"
	"context"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestUserRepositoryDao_FindById(t *testing.T) {
	now := time.Now()
	// 你要去掉毫秒以外的部分
	now = time.UnixMilli(now.UnixMilli())
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)
		ctx  context.Context
		id   int64

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "cache 未命中，查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				d := daomocks.NewMockUserDAO(ctrl)
				d.EXPECT().FindById(gomock.Any(), int64(123)).
					Return(entity.User{
						Id: 123,
						Email: sql.NullString{
							String: "test@gmail.com",
							Valid:  true,
						},
						Password: "test1234",
						Phone: sql.NullString{
							String: "13803711111",
							Valid:  true,
						},
						Ctime: now.UnixMilli(),
						Utime: now.UnixMilli(),
					}, nil)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), gomock.Any()).Return(domain.User{}, redis.ErrKeyNotExist)
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "test@gmail.com",
					Password: "test1234",
					Phone:    "13803711111",
					Ctime:    now,
					Utime:    now,
				}).Return(nil)
				return d, c
			},
			ctx: context.Background(),
			id:  123,

			wantUser: domain.User{
				Id:       123,
				Email:    "test@gmail.com",
				Password: "test1234",
				Phone:    "13803711111",
				Ctime:    now,
				Utime:    now,
			},
			wantErr: nil,
		},
		{
			name: "cache命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				d := daomocks.NewMockUserDAO(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), gomock.Any()).Return(domain.User{
					Id:       123,
					Email:    "test@gmail.com",
					Password: "test1234",
					Phone:    "13803711111",
					Ctime:    now,
					Utime:    now,
				}, nil)
				return d, c
			},
			ctx: context.Background(),
			id:  123,

			wantUser: domain.User{
				Id:       123,
				Email:    "test@gmail.com",
				Password: "test1234",
				Phone:    "13803711111",
				Ctime:    now,
				Utime:    now,
			},
			wantErr: nil,
		},
		{
			name: "cache 未命中，数据库查询失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				d := daomocks.NewMockUserDAO(ctrl)
				d.EXPECT().FindById(gomock.Any(), int64(123)).
					Return(entity.User{}, redis.ErrUnknownForCode)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), gomock.Any()).Return(domain.User{}, redis.ErrKeyNotExist)
				return d, c
			},
			ctx: context.Background(),
			id:  123,

			wantUser: domain.User{},
			wantErr:  redis.ErrUnknownForCode,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ud, uc := tc.mock(ctrl)
			repo := NewUserRepository(ud, uc)
			user, err := repo.FindById(tc.ctx, tc.id)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
			time.Sleep(time.Second)
		})
	}
}
