package repository

import (
	"context"
	"time"

	"github.com/wx-up/go-book/internal/repository/dao/model"

	"github.com/wx-up/go-book/internal/domain"

	"github.com/wx-up/go-book/pkg/logger"

	"github.com/wx-up/go-book/internal/repository/cache"

	"github.com/wx-up/go-book/internal/repository/dao"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizID int64) error

	// IncrLikeCnt 增加点赞数
	IncrLikeCnt(ctx context.Context, biz string, bizID int64, uid int64) error
	// DecrLikeCnt 取消点赞
	DecrLikeCnt(ctx context.Context, biz string, bizID int64, uid int64) error

	// AddCollectionItem 收藏
	// cid 为 0 时表示用户的默认收藏夹
	AddCollectionItem(ctx context.Context, biz string, bizID int64, cid int64, uid int64) error

	// Get  获取某个业务资源的点赞、收藏、阅读
	Get(ctx context.Context, biz string, bizID int64) (domain.Interactive, error)
	// Liked 文章是否被点赞
	Liked(ctx context.Context, biz string, bid int64, uid int64) (bool, error)
	// Collected 文章是否被收藏
	Collected(ctx context.Context, biz string, bid int64, uid int64) (bool, error)
}

type CacheInteractiveRepository struct {
	dao   dao.InteractiveDao
	cache cache.InteractiveCache
	l     logger.Logger
}

func (c *CacheInteractiveRepository) DecrLikeCnt(ctx context.Context, biz string, bizID int64, uid int64) error {
	// TODO implement me
	panic("implement me")
}

func (c *CacheInteractiveRepository) AddCollectionItem(ctx context.Context, biz string, bizID int64, cid int64, uid int64) error {
	// TODO implement me
	panic("implement me")
}

func (c *CacheInteractiveRepository) Get(ctx context.Context, biz string, bizID int64) (domain.Interactive, error) {
	cacheRes, err := c.cache.Get(ctx, biz, bizID)
	if err == nil {
		return cacheRes, nil
	}
	// 查询数据库
	res, err := c.dao.Get(ctx, biz, bizID)
	if err != nil {
		return domain.Interactive{}, err
	}

	inter := c.toDomain(res)

	// 设置缓存
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err = c.cache.Set(ctx, biz, bizID, inter); err != nil {
			// 记录日志
		}
	}()

	return inter, nil
}

func (c *CacheInteractiveRepository) toDomain(res model.Interactive) domain.Interactive {
	return domain.Interactive{}
}

func (c *CacheInteractiveRepository) Liked(ctx context.Context, biz string, bid int64, uid int64) (bool, error) {
	_, err := c.dao.GetLikeInfo(ctx, model.UserLikeBiz{
		Uid:   uid,
		BizId: bid,
		Biz:   biz,
	})
	switch err {
	case nil:
		return true, nil
	case dao.ErrInteractiveLikedNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (c *CacheInteractiveRepository) Collected(ctx context.Context, biz string, bid int64, uid int64) (bool, error) {
	_, err := c.dao.GetCollectionInfo(ctx, model.UserCollectionBiz{
		Uid:   uid,
		BizId: bid,
		Biz:   biz,
	})
	switch err {
	case nil:
		return true, nil
	case dao.ErrInteractiveCollectedNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (c *CacheInteractiveRepository) IncrLikeCnt(ctx context.Context, biz string, bizID int64, uid int64) error {
	err := c.dao.InsertLikeInfo(ctx, biz, bizID, uid)
	if err != nil {
		return err
	}

	// 缓存更新
	if err := c.cache.IncrLikeCntIfPresent(ctx, biz, bizID); err != nil {
		c.l.Error("【缓存】点赞数增加失败", logger.Error(err))
	}
	return nil
}

func (c *CacheInteractiveRepository) IncrReadCnt(ctx context.Context, biz string, bizID int64) error {
	// 更新数据库
	err := c.dao.IncrReadCnt(ctx, biz, bizID)
	if err != nil {
		return err
	}
	// 再更新缓存
	if err := c.cache.IncrReadCntIfPresent(ctx, biz, bizID); err != nil {
		c.l.Error("【缓存】阅读数增加失败", logger.Error(err))
	}
	return nil
}
