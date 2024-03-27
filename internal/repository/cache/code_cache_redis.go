package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

//go:embed lua/send_code.lua
var luaSendCodeScript string

//go:embed lua/verify_code.lua
var luaVerifyCodeScript string

var (
	ErrCodeSendTooMany   = fmt.Errorf("code send too many")
	ErrCodeVerifyTooMany = fmt.Errorf("code verify too many")
	ErrCodeVerifyFail    = fmt.Errorf("code verify fail")
	ErrCodeNotExists     = fmt.Errorf("code not exists")
)

type CodeCache interface {
	Set(ctx context.Context, biz, connect, code string) error // connect 联系方式可以是手机号，也可以是邮箱
	Verify(ctx context.Context, biz, connect, inputCode string) error
}

type RedisCodeCache struct {
	cmd redis.Cmdable
}

func NewRedisCodeCache(cmd redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		cmd: cmd,
	}
}

func (c *RedisCodeCache) Set(ctx context.Context, biz, connect, code string) error {
	key := c.key(biz, connect)
	res, err := c.cmd.Eval(ctx, luaSendCodeScript, []string{key}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		return ErrCodeSendTooMany
	case -2: // 系统错误
		return errors.New("系统错误")
	default:
		return errors.New("系统错误")
	}
}

func (c *RedisCodeCache) Verify(ctx context.Context, biz, connect, inputCode string) error {
	key := c.key(biz, connect)
	res, err := c.cmd.Eval(ctx, luaVerifyCodeScript, []string{key}, inputCode).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		return ErrCodeVerifyTooMany
	case -2:
		return ErrCodeVerifyFail
	case -3:
		return ErrCodeNotExists
	}
	return errors.New("系统错误")
}

func (c *RedisCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("code:%s:%s", biz, phone)
}
