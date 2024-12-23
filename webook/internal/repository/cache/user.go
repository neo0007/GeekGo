package cache

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, user domain.User) error
}

type RedisUserCache struct {
	// redis.Cmdable 是一个接口，传单机 redis 或 cluster redis 等都可以
	client     redis.Cmdable
	expiration time.Duration
}

// A用到了 B，B 一定是接口
// A 用到了 B，B 一定是 A 的字段
// A 用到了 B，A 绝对不初始化 B，而是外面注入

func NewUserCache(client redis.Cmdable) UserCache {
	return &RedisUserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

// err 为 nil user必定有数据
func (cache *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, ErrKeyNotExist
	}
	var user domain.User
	err = json.Unmarshal(val, &user)
	return user, err
}

func (cache *RedisUserCache) Set(ctx context.Context, user domain.User) error {
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := cache.key(user.Id)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()

}

func (cache *RedisUserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
