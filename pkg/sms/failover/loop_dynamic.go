package failover

import (
	"context"
	"sync/atomic"

	"github.com/wx-up/go-book/pkg/sms"
)

// LoopDynamicService 轮训实现：[]sms.Service 负载均衡
type LoopDynamicService struct {
	ss  []sms.Service
	idx uint64
}

func NewLoopDynamicService(ss ...sms.Service) *LoopDynamicService {
	return &LoopDynamicService{ss: ss}
}

func (l *LoopDynamicService) Send(ctx context.Context, tplId string, params []sms.NameArg, phones ...string) error {
	// 为了避免并发问题，这里先将 idx 偏移往后移动一位
	idx := atomic.AddUint64(&l.idx, 1)
	length := uint64(len(l.ss))
	for i := idx; i < length+idx; i++ {
		svc := l.ss[i%length]
		err := svc.Send(ctx, tplId, params, phones...)
		switch err {
		case nil:
			return nil
		case context.Canceled, context.DeadlineExceeded:
			return err
		default:
			// 记录日志
		}
	}
	return ErrAllServicesFailed
}

func (l *LoopDynamicService) Type() string {
	// TODO implement me
	panic("implement me")
}
