package entity

import "database/sql"

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
