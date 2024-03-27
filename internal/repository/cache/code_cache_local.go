package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
)

type (
	CodeItem struct {
		expire time.Time
		cnt    int64
		code   string
	}
	LocalCodeCache struct {
		cache      *lru.Cache[string, *CodeItem]
		lock       sync.RWMutex
		expiration time.Duration
		maps       sync.Map
	}
)

func NewLocalCodeCache(c *lru.Cache[string, *CodeItem], expiration time.Duration) *LocalCodeCache {
	return &LocalCodeCache{
		cache:      c,
		expiration: expiration,
	}
}

func (l *LocalCodeCache) key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func (l *LocalCodeCache) Set(ctx context.Context, biz, connect, code string) error {
	key := l.key(biz, connect)
	now := time.Now()
	l.lock.RLock()
	itm, ok := l.cache.Get(key)
	l.lock.RUnlock()
	if ok {
		if itm.expire.Sub(now) > time.Minute*9 {
			// 不到一分钟
			return ErrCodeSendTooMany
		}
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	_, ok = l.cache.Get(key)
	if ok {
		return ErrCodeSendTooMany
	}
	l.cache.Add(key, &CodeItem{
		code:   code,
		cnt:    3,
		expire: now.Add(l.expiration),
	})
	return nil
}

func (l *LocalCodeCache) SetV2(ctx context.Context, biz, connect, code string) error {
	key := l.key(biz, connect)
	val, _ := l.maps.LoadOrStore(key, &sync.Mutex{})
	lock := val.(*sync.Mutex)
	lock.Lock()
	defer func() {
		l.maps.Delete(key)
		lock.Unlock()
	}()
	// 逻辑
	return nil
}

func (l *LocalCodeCache) Verify(ctx context.Context, biz, connect, inputCode string) error {
	key := l.key(biz, connect)
	l.lock.RLock()
	itm, ok := l.cache.Get(key)
	l.lock.RUnlock()
	if !ok {
		return ErrCodeNotExists
	}
	if itm.cnt <= 0 {
		return ErrCodeVerifyTooMany
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	itm.cnt--
	if itm.code != inputCode {
		return ErrCodeVerifyFail
	}
	return nil
}
