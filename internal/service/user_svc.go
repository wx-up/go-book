package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/wx-up/go-book/internal/repository"

	"github.com/wx-up/go-book/internal/domain"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
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
