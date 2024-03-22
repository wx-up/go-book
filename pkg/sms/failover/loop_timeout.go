package failover

import (
	"context"
	"sync/atomic"

	"github.com/wx-up/go-book/pkg/sms"
)

// LoopTimeoutService 连续N个超时响应就切换
type LoopTimeoutService struct {
	ss []sms.Service

	idx int32

	// 连续超时的个数
	cnt int32

	// 超时个数阈值
	threshold int32
}

func NewLoopTimeoutService(cnt int32, threshold int32, ss ...sms.Service) *LoopTimeoutService {
	return &LoopTimeoutService{
		ss:  ss,
		idx: 0,
		cnt: cnt,
	}
}

func (l *LoopTimeoutService) Send(ctx context.Context, tplId string, params []sms.NameArg, phones ...string) error {
	idx := atomic.LoadInt32(&l.idx)
	cnt := atomic.LoadInt32(&l.cnt)
	if cnt >= l.threshold {
		// 取余数是为了避免索引越界
		newIdx := (idx + 1) % int32(len(l.ss))
		// CAS 操作失败，说明并发了，其他人改了
		if atomic.CompareAndSwapInt32(&l.idx, idx, newIdx) {
			atomic.StoreInt32(&l.cnt, 0)
		}
		// 重新获取索引
		idx = atomic.LoadInt32(&l.idx)
	}
	err := l.ss[idx].Send(ctx, tplId, params, phones...)
	switch err {
	case nil:
		// 没有错误，重制次数
		atomic.StoreInt32(&l.cnt, 0)
	case context.DeadlineExceeded:
		// 超时，计数加1
		atomic.AddInt32(&l.cnt, 1)

	default:
		// 其他异常，这里考虑可以切换到下一个服务
		// 暂时什么也不做
	}
	return err
}

func (l *LoopTimeoutService) Type() string {
	// TODO implement me
	panic("implement me")
}
