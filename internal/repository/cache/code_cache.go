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

type CodeCache struct {
	cmd redis.Cmdable
}

func NewCodeCache(cmd redis.Cmdable) *CodeCache {
	return &CodeCache{
		cmd: cmd,
	}
}

func (c *CodeCache) Set(ctx context.Context, biz, phone, code string) error {
	key := c.key(biz, phone)
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

func (c *CodeCache) Verify(ctx context.Context, biz, phone, inputCode string) error {
	key := c.key(biz, phone)
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

func (c *CodeCache) key(biz, phone string) string {
	return fmt.Sprintf("code:%s:%s", biz, phone)
}
