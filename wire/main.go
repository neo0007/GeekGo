package wire

import (
	"Neo/Workplace/goland/src/GeekGo/wire/repository"
	"Neo/Workplace/goland/src/GeekGo/wire/repository/dao"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		panic(err)
	}
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	fmt.Printf("v%v", repo)
}
