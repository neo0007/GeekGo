package gorm

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository/dao/entity"
	"gorm.io/gorm"
)

func InitTable(db *gorm.DB) error {
	//如果没有建好表，或者表需要修改，则 return 下面语句
	return db.AutoMigrate(&entity.User{})
}
