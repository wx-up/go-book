package repository

import (
	"context"

	"github.com/wx-up/go-book/internal/repository/cache"
)

var (
	ErrCodeSendTooMany   = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany
)

type CodeRepository struct {
	c *cache.CodeCache
}

func NewCodeRepository(c *cache.CodeCache) *CodeRepository {
	return &CodeRepository{c: c}
}

func (cr *CodeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return cr.c.Set(ctx, biz, phone, code)
}

func (cr *CodeRepository) Verify(ctx context.Context, biz, phone, code string) error {
	return cr.c.Verify(ctx, biz, phone, code)
}
