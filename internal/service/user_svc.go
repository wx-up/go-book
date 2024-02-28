package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/wx-up/go-book/internal/repository"

	"github.com/wx-up/go-book/internal/domain"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("账号或者密码不对")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) Login(ctx context.Context, obj domain.User) (domain.User, error) {
	// 查找用户
	u, err := svc.repo.FindByEmail(ctx, obj.Email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}

	// 比较密码
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(obj.Password)); err != nil {
		// 错误被转换了，需要打日志
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *UserService) SignUp(ctx context.Context, obj domain.User) error {
	// 密码加密
	hash, err := bcrypt.GenerateFromPassword([]byte(obj.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	obj.Password = string(hash)

	return svc.repo.Create(ctx, obj)
}
