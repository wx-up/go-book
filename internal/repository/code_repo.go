package repository

import (
	"context"

	"github.com/wx-up/go-book/internal/repository/cache"
)

var (
	ErrCodeSendTooMany   = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany
)

type CodeRepository interface {
	Store(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) error
}

type CacheCodeRepository struct {
	c cache.CodeCache
}

func NewCacheCodeRepository(c cache.CodeCache) CodeRepository {
	return &CacheCodeRepository{c: c}
}

func (cr *CacheCodeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return cr.c.Set(ctx, biz, phone, code)
}

func (cr *CacheCodeRepository) Verify(ctx context.Context, biz, phone, code string) error {
	return cr.c.Verify(ctx, biz, phone, code)
}
