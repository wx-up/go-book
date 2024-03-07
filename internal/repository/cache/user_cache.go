package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wx-up/go-book/internal/domain"
)

type UserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func NewUserCache(cmd redis.Cmdable) *UserCache {
	return &UserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}

func (uc *UserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := uc.key(u.Id)
	return uc.cmd.Set(ctx, key, val, uc.expiration).Err()
}

var ErrKeyNotExist = fmt.Errorf("key not exist")

// Get 需要区分是 redis报错 还是 key 不存在
func (uc *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := uc.key(id)
	val, err := uc.cmd.Get(ctx, key).Bytes()
	// key 不存在的时，redis 报错 redis.Nil
	if err != nil {
		if err == redis.Nil {
			return domain.User{}, ErrKeyNotExist
		}
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(val, &u)
	return u, err
}

func (uc *UserCache) key(id int64) string {
	const keyPattern = "user:info:%d"
	return fmt.Sprintf(keyPattern, id)
}
