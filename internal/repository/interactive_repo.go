package repository

import (
	"context"

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
}

type CacheInteractiveRepository struct {
	dao   dao.InteractiveDao
	cache cache.InteractiveCache
	l     logger.Logger
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
