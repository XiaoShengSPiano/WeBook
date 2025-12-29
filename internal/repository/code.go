package repository

import (
	"context"
	"webook/internal/repository/cache"
)

var (
	ErrCodeSendTooFrequently  = cache.ErrCodeSendTooFrequently
	ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
)

type CodeRepository interface {
	Store(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) error
}

type CacheCodeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(c cache.CodeCache) CodeRepository {
	return &CacheCodeRepository{
		cache: c,
	}
}

func (r *CacheCodeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return r.cache.SetCode(ctx, biz, phone, code)
}

func (r *CacheCodeRepository) Verify(ctx context.Context, biz, phone, inputCode string) error {
	return r.cache.VerifyCode(ctx, biz, phone, inputCode)
}
