package repository

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFound = dao.ErrUserNotFound

// 存储层
type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: c,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		// 缓存中必定存在数据
		return u, nil
	}

	if err == cache.ErrkeyNotExists {
		// 缓存中不存在数据，查询数据库
	}

	// 缓存中不存在数据，查询数据库
	ur, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	u = domain.User{
		Id:       ur.Id,
		Email:    ur.Email,
		Password: ur.Password,
	}
	// 缓存写回
	err = r.cache.Set(ctx, u)
	if err != nil {
		// 缓存写入失败，打日志，做监控
	}
	return u, err

	// 存在别的错误（可能redis崩溃，需要保护数据库）
}
