package dao

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository/dao/entity"
	"context"
)

type UserDAO interface {
	Insert(c context.Context, u entity.User) error
	FindByEmail(ctx context.Context, email string) (entity.User, error)
	FindByPhone(ctx context.Context, phone string) (entity.User, error)
	FindById(ctx context.Context, id int64) (entity.User, error)
	FindByWechat(ctx context.Context, openID string) (entity.User, error)
}
