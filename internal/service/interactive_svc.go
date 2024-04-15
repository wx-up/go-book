package service

import (
	"context"

	"github.com/wx-up/go-book/internal/repository"
)

// InteractiveService 阅读服务
type InteractiveService interface {
	// IncrReadCount 阅读计数
	// biz 业务类型
	// bid 业务ID
	IncrReadCount(ctx context.Context, biz string, bid int64) error

	// Like 点赞
	Like(ctx context.Context, biz string, bid int64, uid int64) error
	// CancelLike 取消点赞
	CancelLike(ctx context.Context, biz string, bid int64, uid int64) error
}

type interactiveService struct {
	repo repository.InteractiveRepository
}

func (i *interactiveService) Like(ctx context.Context, biz string, bid int64, uid int64) error {
	// TODO implement me
	panic("implement me")
}

func (i *interactiveService) CancelLike(ctx context.Context, biz string, bid int64, uid int64) error {
	// TODO implement me
	panic("implement me")
}

func (i *interactiveService) IncrReadCount(ctx context.Context, biz string, bid int64) error {
	return i.repo.IncrReadCnt(ctx, biz, bid)
}
