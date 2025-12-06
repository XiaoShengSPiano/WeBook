package repository

import (
	"context"
	"webook/internal/repository/cache"
)

var (
	ErrCodeSendTooFrequently  = cache.ErrCodeSendTooFrequently
	ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
)

type CodeRepository struct {
	cache *cache.CodeCache
}

func NewCodeRepository(c *cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		cache: c,
	}
}

func (r *CodeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return r.cache.SetCode(ctx, biz, phone, code)
}

func (r *CodeRepository) Verify(ctx context.Context, biz, phone, inputCode string) error {
	return r.cache.VerifyCode(ctx, biz, phone, inputCode)
}
