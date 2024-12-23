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
	ErrUserDuplicate = errors.New("邮箱或手机号码冲突了！")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Insert(c context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
}

type GORMUserDAO struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}

func (dao *GORMUserDAO) Insert(c context.Context, u User) error {
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

func (dao *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).First(&u, "email = ?", email).Error
	//err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).First(&u, "phone = ?", phone).Error
	//err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).First(&u, "id = ?", id).Error
	return u, err
}

// User 对应数据库表结构
// 有些人叫做 entity，有些人叫做 model，有些人叫做 PO(Persistent Object)
type User struct {
	Id       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"type:varchar(100);unique"`
	Password string         `gorm:"type:varchar(100)"`
	// 唯一索引允许有多个空值
	// 但是不能有多个空字符串“”
	Phone    sql.NullString `gorm:"type:varchar(100);unique"`
	Nickname string         `gorm:"type:varchar(100)"`
	Birthday string         `gorm:"type:varchar(100)"`
	AboutMe  string         `gorm:"type:varchar(255)"`
	//创建时间, 毫秒数
	Ctime int64
	//更新时间，毫秒数
	Utime int64
}
