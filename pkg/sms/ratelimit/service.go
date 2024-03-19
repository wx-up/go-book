package ratelimit

import (
	"context"
	"fmt"

	"github.com/wx-up/go-book/pkg/ratelimit"
	"github.com/wx-up/go-book/pkg/sms"
)

var ErrLimited = fmt.Errorf("短信服务限流，请稍后再试")

type Service struct {
	sms     sms.Service
	limiter ratelimit.Limiter
}

func NewService(sms sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &Service{
		sms:     sms,
		limiter: limiter,
	}
}

func (s *Service) Send(ctx context.Context, tplId string, params []sms.NameArg, phones ...string) error {
	limited, err := s.limiter.Limit(ctx, fmt.Sprintf("sms:%s", s.sms.Type()))
	if err != nil {
		// 这里一般是因为 redis 蹦了，有两种处理策略
		// 限流：保守策略，比如你对接的短信服务商技术能力比较弱的时候
		// 不限流，直接放行：比如你的短信服务商技术能力很强，本身就有限流策略
		// 或者你的业务可用性要求很高，那就采用该策略。
		// 这里采用保守策略
		return fmt.Errorf("短信服务限流出现问题：%w", err)
	}
	if limited {
		return ErrLimited
	}
	return s.sms.Send(ctx, tplId, params, phones...)
}

func (s *Service) Type() string {
	return "sms-ratelimit"
}
