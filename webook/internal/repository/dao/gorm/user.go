package gorm

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository/dao"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository/dao/entity"
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicate = errors.New("邮箱或手机号码冲突了！")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

type GORMUserDAO struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) dao.UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}

func (dao *GORMUserDAO) Insert(c context.Context, u entity.User) error {
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	err := dao.db.WithContext(c).Create(&u).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			//邮箱冲突或手机号码冲突
			return ErrUserDuplicate
		}
	}
	return err
}

func (dao *GORMUserDAO) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	var u entity.User
	err := dao.db.WithContext(ctx).First(&u, "email = ?", email).Error
	//err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (entity.User, error) {
	var u entity.User
	err := dao.db.WithContext(ctx).First(&u, "phone = ?", phone).Error
	//err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) FindById(ctx context.Context, id int64) (entity.User, error) {
	var u entity.User
	err := dao.db.WithContext(ctx).First(&u, "id = ?", id).Error
	return u, err
}

func (dao *GORMUserDAO) FindByWechat(ctx context.Context, openID string) (entity.User, error) {
	var u entity.User
	err := dao.db.WithContext(ctx).First(&u, "wechat_open_id = ?", openID).Error
	//err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}
