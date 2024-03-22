package failover

import (
	"context"
	"errors"
	"log"

	"github.com/wx-up/go-book/pkg/sms"
)

var ErrAllServicesFailed = errors.New("all services failed")

// LoopService 轮训实现
type LoopService struct {
	ss []sms.Service
}

func NewLoopService(ss ...sms.Service) *LoopService {
	return &LoopService{ss: ss}
}

func (l *LoopService) Send(ctx context.Context, tplId string, params []sms.NameArg, phones ...string) error {
	for _, s := range l.ss {
		err := s.Send(ctx, tplId, params, phones...)
		if err == nil {
			return nil
		}
		if err == context.DeadlineExceeded {
			return errors.New("context timeout")
		}
		// 记录日志
		log.Println(err)
	}
	return ErrAllServicesFailed
}

func (l *LoopService) Type() string {
	// TODO implement me
	panic("implement me")
}
