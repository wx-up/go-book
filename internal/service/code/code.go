package code

import (
	"context"
)

// Service 验证码服务
type Service interface {
	Send(ctx context.Context, biz string, phone string) error
	// Verify 验证码校验
	// error 表示系统错误 bool 表示验证是否通过
	// 或者也可以设计成只返回 error 系统错误直接返回 error 验证不通过返回一个特殊的错误
	Verify(ctx context.Context, biz string, phone string, code string) (bool, error)
}
