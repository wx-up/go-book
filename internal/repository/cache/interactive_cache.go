package cache

import (
	"context"
	_ "embed"
	"fmt"
	"strconv"
	"time"

	"github.com/wx-up/go-book/internal/domain"

	"github.com/redis/go-redis/v9"
)

//go:embed lua/interactive_incr_cnt.lua
var luaIncrCnt string

var InteractiveKeyNotExits = fmt.Errorf("interactive key not exists")

const (
	fieldCollectCnt = "collect_cnt"
	fieldReadCnt    = "read_cnt"
	fieldLikeCnt    = "like_cnt"
)

type InteractiveCache interface {
	IncrCollectCntIfPresent(ctx context.Context, biz string, bid int64) error
	IncrReadCntIfPresent(ctx context.Context, biz string, bid int64) error
	IncrLikeCntIfPresent(ctx context.Context, biz string, bid int64) error
	// CancelLikeCntIfPresent 取消点赞
	CancelLikeCntIfPresent(ctx context.Context, biz string, bid int64) error

	Get(ctx context.Context, biz string, bid int64) (domain.Interactive, error)
	Set(ctx context.Context, biz string, bid int64, inter domain.Interactive) error
}

type RedisInteractiveCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func (r *RedisInteractiveCache) Set(ctx context.Context, biz string, bid int64, inter domain.Interactive) error {
	key := r.key(biz, bid)
	err := r.client.HSet(ctx, key, fieldCollectCnt, inter.CollectCnt,
		fieldReadCnt, inter.ReadCnt,
		fieldLikeCnt, inter.LikeCnt,
	).Err()
	if err != nil {
		return err
	}
	return r.client.Expire(ctx, key, r.expiration).Err()
}

func (r *RedisInteractiveCache) Get(ctx context.Context, biz string, bid int64) (domain.Interactive, error) {
	res, err := r.client.HGetAll(ctx, r.key(biz, bid)).Result()
	if err != nil {
		return domain.Interactive{}, err
	}
	if len(res) == 0 {
		return domain.Interactive{}, InteractiveKeyNotExits
	}
	var inter domain.Interactive
	// 这边是可以忽略错误的
	inter.CollectCnt, _ = strconv.ParseInt(res[fieldCollectCnt], 10, 64)
	inter.LikeCnt, _ = strconv.ParseInt(res[fieldLikeCnt], 10, 64)
	inter.ReadCnt, _ = strconv.ParseInt(res[fieldReadCnt], 10, 64)
	return inter, nil
}

func (r *RedisInteractiveCache) CancelLikeCntIfPresent(ctx context.Context, biz string, bid int64) error {
	return r.wrapRes(func() (int64, error) {
		return r.client.Eval(ctx, luaIncrCnt, []string{r.key(biz, bid)}, fieldLikeCnt, -1).Int64()
	})
}

func NewRedisInteractiveCache(client redis.Cmdable, expiration time.Duration) *RedisInteractiveCache {
	return &RedisInteractiveCache{
		client:     client,
		expiration: expiration,
	}
}

func (r *RedisInteractiveCache) IncrCollectCntIfPresent(ctx context.Context, biz string, bid int64) error {
	return r.wrapRes(func() (int64, error) {
		return r.client.Eval(ctx, luaIncrCnt, []string{r.key(biz, bid)}, fieldCollectCnt, 1).Int64()
	})
}

func (r *RedisInteractiveCache) wrapRes(handle func() (int64, error)) error {
	res, err := handle()
	if err != nil {
		return err
	}
	switch res {
	case 1:
		return nil
	case 0:
		return InteractiveKeyNotExits
	default:
		return fmt.Errorf("interactive unexpected result: %d", res)
	}
}

func (r *RedisInteractiveCache) key(biz string, bid int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bid)
}

func (r *RedisInteractiveCache) IncrReadCntIfPresent(ctx context.Context, biz string, bid int64) error {
	return r.wrapRes(func() (int64, error) {
		return r.client.Eval(ctx, luaIncrCnt, []string{r.key(biz, bid)}, fieldReadCnt, 1).Int64()
	})
}

func (r *RedisInteractiveCache) IncrLikeCntIfPresent(ctx context.Context, biz string, bid int64) error {
	return r.wrapRes(func() (int64, error) {
		return r.client.Eval(ctx, luaIncrCnt, []string{r.key(biz, bid)}, fieldLikeCnt, 1).Int64()
	})
}
