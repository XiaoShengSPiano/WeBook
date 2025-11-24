package service

import (
	"context"
	"errors"
	"webook/internal/domain"
	"webook/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
var ErrInvalidUserOrPassword = errors.New("账号/密码不对......")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	// 对密码进行加密
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hash)

	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, u domain.User) (domain.User, error) {
	// 查询用户是否存在
	ur, err := svc.repo.FindByEmail(ctx, u.Email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	if err != nil {
		return domain.User{}, err
	}

	// 验证密码是否正确
	err = bcrypt.CompareHashAndPassword([]byte(ur.Password), []byte(u.Password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	return ur, nil
}

func (svc *UserService) Profile(ctx context.Context, userId int64) (domain.User, error) {
	// TODO
	u, err := svc.repo.FindById(ctx, userId)

	return u, err
}
