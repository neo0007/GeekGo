package ioc

import (
	"Neo/Workplace/goland/src/GeekGo/webook/config"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		// 应该只在初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		// 一旦初始化出错，即退出
		panic(err)
	}
	initDB := true
	if initDB {
		err = dao.InitTable(db)
		if err != nil {
			panic(err)
		}
	}
	return db
}
