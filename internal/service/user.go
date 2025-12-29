package service

import (
	"context"
	"errors"
	"webook/internal/domain"
	"webook/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicate
var ErrInvalidUserOrPassword = errors.New("账号/密码不对......")

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, u domain.User) (domain.User, error)
	Profile(ctx context.Context, userId int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (svc *userService) SignUp(ctx context.Context, u domain.User) error {
	// 对密码进行加密
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hash)

	return svc.repo.Create(ctx, u)
}

func (svc *userService) Login(ctx context.Context, u domain.User) (domain.User, error) {
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

func (svc *userService) Profile(ctx context.Context, userId int64) (domain.User, error) {
	// TODO
	u, err := svc.repo.FindById(ctx, userId)

	return u, err
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	// 查询用户是否存在(快路径)
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		return u, err
	}

	// 如果用户不存在，则创建新用户（插入数据，满路经）
	u = domain.User{
		Phone: phone,
	}
	err = svc.repo.Create(ctx, u)
	// 系统错误
	if err != nil && err != repository.ErrUserDuplicate {
		return u, err
	}

	return svc.repo.FindByPhone(ctx, phone)
}
