package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (dao *UserDao) Insert(c context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	return dao.db.WithContext(c).Create(&u).Error
}

// User 对应数据库表结构
// 有些人叫做 entity，有些人叫做 model，有些人叫做 PO(Persistent Object)
type User struct {
	Id       int64  `gorm:"primary_key,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string

	//创建时间, 毫秒数
	Ctime int64
	//更新时间，毫秒数
	Utime int64
}
