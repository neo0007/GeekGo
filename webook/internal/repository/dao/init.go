package dao

import "gorm.io/gorm"

func InitTable(db *gorm.DB) error {
	//如果没有建好表，或者表需要修改，则 return 下面语句
	return db.AutoMigrate(&User{})
}
