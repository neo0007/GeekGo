package repository

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/domain"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository/cache"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository/dao"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository/dao/entity"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository/dao/gorm"
	"context"
	"database/sql"
	"log"
	"time"
)

var ErrUserDuplicate = gorm.ErrUserDuplicate
var ErrUserNotFound = gorm.ErrUserNotFound

type UserRepository interface {
	Create(c context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByWechat(ctx context.Context, openID string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
}

func (r *UserRepositoryDao) FindByWechat(ctx context.Context, openID string) (domain.User, error) {
	ud, err := r.dao.FindByWechat(ctx, openID)
	if err != nil {
		return domain.User{}, err
	}
	u := r.entityToDomain(ud)
	err = r.cache.Set(ctx, u)
	if err != nil {
		log.Println(err)
	}
	return u, err

}

type UserRepositoryDao struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDAO, c cache.UserCache) UserRepository {
	return &UserRepositoryDao{
		dao:   dao,
		cache: c,
	}
}

func (r *UserRepositoryDao) Create(c context.Context, u domain.User) error {
	return r.dao.Insert(c, r.domainToEntity(u))
}

func (r *UserRepositoryDao) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	ud, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(ud), nil

	//先从 cache 里面找
	//再从 dao 里面找
	//找到了回写 cache
}

func (r *UserRepositoryDao) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	ud, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	u := r.entityToDomain(ud)
	err = r.cache.Set(ctx, u)
	if err != nil {
		log.Println(err)
	}
	return u, err

	//先从 cache 里面找
	//再从 dao 里面找
	//找到了回写 cache
}

func (r *UserRepositoryDao) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		// 必然有数据
		return u, nil
	}
	//没有这个数据
	//if err == cache.ErrKeyNotExist {
	////	 去数据库里加载
	//}
	ud, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = r.entityToDomain(ud)

	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			log.Println(err)
		}
	}()

	return u, err
}

func (r *UserRepositoryDao) entityToDomain(u entity.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Phone:    u.Phone.String,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		AboutMe:  u.AboutMe,
		WechatInfo: domain.WechatInfo{
			OpenID:  u.WechatOpenID.String,
			UnionID: u.WechatUnionID.String,
		},
		Ctime: time.UnixMilli(u.Ctime),
		Utime: time.UnixMilli(u.Utime),
	}
}

func (r *UserRepositoryDao) domainToEntity(u domain.User) entity.User {
	return entity.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Password: u.Password,
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		AboutMe:  u.AboutMe,
		WechatOpenID: sql.NullString{
			String: u.WechatInfo.OpenID,
			Valid:  u.WechatInfo.OpenID != "",
		},
		WechatUnionID: sql.NullString{
			String: u.WechatInfo.UnionID,
			Valid:  u.WechatInfo.UnionID != "",
		},
		Ctime: u.Ctime.UnixMilli(),
		Utime: u.Utime.UnixMilli(),
	}
}
