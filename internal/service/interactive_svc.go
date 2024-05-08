package service

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/wx-up/go-book/internal/domain"

	"github.com/wx-up/go-book/internal/repository"
)

// InteractiveService 互动服务
//
//go:generate mockgen -destination=./mocks/interactive_svc.mock.go -package=svcmocks -source=interactive_svc.go InteractiveService
type InteractiveService interface {
	// IncrReadCount 阅读计数
	// biz 业务类型
	// bid 业务ID
	IncrReadCount(ctx context.Context, biz string, bid int64) error

	// Like 点赞
	Like(ctx context.Context, biz string, bid int64, uid int64) error
	// CancelLike 取消点赞
	CancelLike(ctx context.Context, biz string, bid int64, uid int64) error

	Get(ctx context.Context, biz string, bid int64, uid int64) (domain.Interactive, error)

	Liked(ctx context.Context, biz string, bid int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, bid int64, uid int64) (bool, error)

	GetByIds(ctx context.Context, biz string, ids []int64) ([]domain.Interactive, error)
}

type interactiveService struct {
	repo repository.InteractiveRepository
}

func (i *interactiveService) GetByIds(ctx context.Context, biz string, ids []int64) ([]domain.Interactive, error) {
	// TODO implement me
	panic("implement me")
}

func NewInteractiveService(repo repository.InteractiveRepository) InteractiveService {
	return &interactiveService{
		repo: repo,
	}
}

func (i *interactiveService) Liked(ctx context.Context, biz string, bid int64, uid int64) (bool, error) {
	return i.repo.Liked(ctx, biz, bid, uid)
}

func (i *interactiveService) Collected(ctx context.Context, biz string, bid int64, uid int64) (bool, error) {
	return i.repo.Collected(ctx, biz, bid, uid)
}

func (i *interactiveService) Get(ctx context.Context, biz string, bid int64, uid int64) (domain.Interactive, error) {
	// 文章的点赞等指标进行了缓存
	// 某个用户对某个文章是否点赞和收藏没有必要缓存，因为实际场景下，不会有用户不停刷新一篇文章的，
	var (
		inter     domain.Interactive
		liked     bool
		collected bool
		eg        errgroup.Group
	)
	eg.Go(func() error {
		var err error
		inter, err = i.repo.Get(ctx, biz, bid)
		return err
	})

	// 用户针对该资源是否点赞、收藏
	eg.Go(func() error {
		var err error
		liked, err = i.Liked(ctx, biz, bid, uid)
		return err
	})
	eg.Go(func() error {
		var err error
		collected, err = i.Collected(ctx, biz, bid, uid)
		return err
	})
	if err := eg.Wait(); err != nil {
		return domain.Interactive{}, err
	}
	inter.Liked = liked
	inter.Collected = collected
	return inter, nil
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
