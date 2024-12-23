//go:build wireinject

package wire

import (
	"Neo/Workplace/goland/src/GeekGo/wire/repository"
	"Neo/Workplace/goland/src/GeekGo/wire/repository/dao"
	"github.com/google/wire"
)

func InitRepository() *repository.UserRepository {
	wire.Build(repository.NewUserRepository, dao.NewUserDAO, InitDB)
	return new(repository.UserRepository)
}
