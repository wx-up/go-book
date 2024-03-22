package retryable

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/wx-up/go-book/pkg/sms"
)

type Service struct {
	sms      sms.Service
	retryMax int32
	cnt      int32
}

func (s *Service) Send(ctx context.Context, tplId string, params []sms.NameArg, phones ...string) error {
	err := s.sms.Send(ctx, tplId, params, phones...)
	if err == nil {
		return nil
	}
	atomic.AddInt32(&s.cnt, 1)
	for s.cnt <= s.retryMax {
		err = s.sms.Send(ctx, tplId, params, phones...)
		if err == nil {
			return nil
		}
		atomic.AddInt32(&s.cnt, 1)
	}
	return errors.New("重试失败了")
}

func (s *Service) Type() string {
	// TODO implement me
	panic("implement me")
}
