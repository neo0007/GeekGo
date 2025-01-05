package repository

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository/cache"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository/cache/redis"
	"context"
)

var (
	ErrCodeSendTooMany        = redis.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = redis.ErrCodeVerifyTooManyTimes
	ErrUnknownForCode         = redis.ErrUnknownForCode
)

type CodeRepository interface {
	Store(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type CodeRepositoryCache struct {
	cache cache.CodeCache
}

func NewCodeRepository(cache cache.CodeCache) CodeRepository {
	return &CodeRepositoryCache{
		cache: cache,
	}
}

func (repo *CodeRepositoryCache) Store(ctx context.Context, biz string,
	phone string, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}

func (repo *CodeRepositoryCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return repo.cache.Verify(ctx, biz, phone, inputCode)
}
