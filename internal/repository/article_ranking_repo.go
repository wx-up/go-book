package repository

import (
	"context"

	"github.com/wx-up/go-book/internal/repository/cache"

	"github.com/wx-up/go-book/internal/domain"
)

type ArticleRankingRepo interface {
	Set(ctx context.Context, objs []domain.Article) error
	Get(ctx context.Context) ([]domain.Article, error)
}

type CacheArticleRankingRepo struct {
	redisCache cache.ArticleRankingCache
	localCache cache.ArticleRankingCache
}

func NewCacheArticleRankingRepo(
	cache cache.ArticleRankingCache,
	localCache cache.ArticleRankingCache,
) *CacheArticleRankingRepo {
	return &CacheArticleRankingRepo{
		redisCache: cache,
		localCache: localCache,
	}
}

func (c *CacheArticleRankingRepo) Get(ctx context.Context) ([]domain.Article, error) {
	data, err := c.localCache.Get(ctx)
	if err == nil {
		return data, nil
	}
	data, err = c.redisCache.Get(ctx)
	if err == nil {
		_ = c.localCache.Set(ctx, data)
	} else {
		return c.localCache.ForceGet(ctx)
	}
	return data, err
}

func (c *CacheArticleRankingRepo) Set(ctx context.Context, objs []domain.Article) error {
	// 操作本地缓存基本上不会出错
	_ = c.localCache.Set(ctx, objs)
	return c.redisCache.Set(ctx, objs)
}
