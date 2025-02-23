package local

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository/cache"
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository/cache/redis"
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/golang-lru"
	"sync"
	"time"
)

var (
	ErrCodeSendTooMany        = errors.New("发送验证码太频繁")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数太多")
	ErrUnknownForCode         = errors.New("code未知错误")
)

type LocalCodeCache struct {
	cache      *lru.Cache
	lock       sync.Mutex
	expiration time.Duration
	maps       sync.Map
}

func NewLocalCodeCache(c *lru.Cache, expiration time.Duration) cache.CodeCache {
	return &LocalCodeCache{
		cache:      c,
		expiration: expiration,
	}
}

func (l *LocalCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	key := l.key(biz, phone)
	now := time.Now()
	val, ok := l.cache.Get(key)
	if !ok {
		//说明没有验证码
		l.cache.Add(key, codeItem{
			code:   code,
			cnt:    3,
			expire: now.Add(l.expiration),
		})
		return nil
	}
	itm, ok := val.(codeItem)
	if !ok {
		return errors.New("系统错误")
	}
	if itm.expire.Sub(now) > time.Minute*9 {
		return ErrCodeSendTooMany
	}
	//重发
	l.cache.Add(key, codeItem{
		code:   code,
		cnt:    3,
		expire: now.Add(l.expiration),
	})
	return nil
}

func (l *LocalCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	//方法一：
	//l.lock.Lock()
	//defer l.lock.Unlock()

	key := l.key(biz, phone)

	//方法二：
	//如果的你的 key 非常多，maps 本身就占据了很多内存
	//lock, _ := l.maps.LoadOrStore(key, &sync.Map{})
	//lock.(*sync.Mutex).Lock()
	//defer lock.(*sync.Mutex).Unlock()

	//方法三：
	lock, _ := l.maps.LoadOrStore(key, &sync.Map{})
	lock.(*sync.Mutex).Lock()
	defer func() {
		l.maps.Delete(key)
		lock.(*sync.Mutex).Unlock()
	}()

	val, ok := l.cache.Get(key)
	if !ok {
		return false, redis.ErrKeyNotExist
	}
	itm, ok := val.(codeItem)
	if !ok {
		return false, errors.New("系统错误")
	}
	if itm.cnt <= 0 {
		return false, ErrCodeVerifyTooManyTimes
	}
	itm.cnt--
	return itm.code == inputCode, nil
}

func (c *LocalCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

type codeItem struct {
	code   string
	cnt    int
	expire time.Time
}
