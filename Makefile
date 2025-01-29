.PHONY: mock
mock:
	@mockgen -source=webook/internal/service/user.go -package=svcmocks -destination=webook/internal/service/mocks/user.mock.go
	@mockgen -source=webook/internal/service/code.go -package=svcmocks -destination=webook/internal/service/mocks/code.mock.go
	@mockgen -source=webook/internal/repository/user.go -package=repomocks -destination=webook/internal/repository/mocks/user.mock.go
	@mockgen -source=webook/internal/repository/code.go -package=repomocks -destination=webook/internal/repository/mocks/code.mock.go
	@mockgen -source=webook/internal/repository/dao/types.go -package=daomocks -destination=webook/internal/repository/dao/mocks/types.mock.go
	@mockgen -source=webook/internal/repository/cache/types.go -package=cachemocks -destination=webook/internal/repository/cache/mocks/types.mock.go
	@go mod tidy
