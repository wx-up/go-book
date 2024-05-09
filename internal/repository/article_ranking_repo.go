package repository

import (
	"context"

	"github.com/wx-up/go-book/internal/repository/cache"

	"github.com/wx-up/go-book/internal/domain"
)

type ArticleRankingRepo interface {
	Set(ctx context.Context, objs []domain.Article) error
}

type CacheArticleRankingRepo struct {
	cache cache.ArticleRankingCache
}

func (c *CacheArticleRankingRepo) Set(ctx context.Context, objs []domain.Article) error {
	return c.cache.Set(ctx, objs)
}

func NewCacheArticleRankingRepo(cache cache.ArticleRankingCache) *CacheArticleRankingRepo {
	return &CacheArticleRankingRepo{cache: cache}
}
