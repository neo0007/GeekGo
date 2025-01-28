.PHONY: mock
mock:
	@mockgen -source=webook/internal/service/user.go -package=svcmocks -destination=webook/internal/service/mocks/user.mock.go
	@mockgen -source=webook/internal/service/code.go -package=svcmocks -destination=webook/internal/service/mocks/code.mock.go
	@mockgen -source=webook/internal/repository/user.go -package=repomocks -destination=webook/internal/repository/mocks/user.mock.go
	@mockgen -source=webook/internal/repository/code.go -package=repomocks -destination=webook/internal/repository/mocks/code.mock.go
	@mockgen -source=webook/internal/repository/dao/gorm/user.go -package=gormmocks -destination=webook/internal/repository/dao/mocks/user.mock.go
	@mockgen -source=webook/internal/repository/cache/redis/code.go -package=redismocks -destination=webook/internal/repository/cache/redis/mocks/code.mock.go
	@mockgen -source=webook/internal/repository/cache/redis/user.go -package=redismocks -destination=webook/internal/repository/cache/redis/mocks/user.mock.go
	@go mod tidy
