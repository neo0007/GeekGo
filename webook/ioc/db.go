package ioc

import (
	myGorm "Neo/Workplace/goland/src/GeekGo/webook/internal/repository/dao/gorm"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	//dsn := viper.GetString("db.mysql.dsn")
	//println("dsn:", dsn)
	//db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	//设置默认值
	cfg := Config{
		DSN: "root:root@tcp(localhost:13316)/webook_default",
	}

	err := viper.UnmarshalKey("db.mysql", &cfg)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.Open(cfg.DSN))
	if err != nil {
		// 应该只在初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		// 一旦初始化出错，即退出
		panic(err)
	}

	err = myGorm.InitTable(db)
	if err != nil {
		panic(err)
	}

	return db
}
