package service

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"golang.org/x/crypto/bcrypt"

	"github.com/wx-up/go-book/internal/repository"

	"github.com/wx-up/go-book/internal/domain"
)

var (
	ErrUserDuplicate         = repository.ErrUserDuplicate
	ErrInvalidUserOrPassword = errors.New("账号或者密码不对")
)

type UserService interface {
	Login(ctx context.Context, obj domain.User) (domain.User, error)
	SignUp(ctx context.Context, obj domain.User) error
	FindOrCreateByPhone(ctx context.Context, phone string) (domain.User, error)
	Profile(ctx context.Context, uid int64) (domain.User, error)
	FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error)
}

type userService struct {
	repo   repository.UserRepository
	logger *zap.Logger
}

func NewUserService(repo repository.UserRepository, logger *zap.Logger) UserService {
	return &userService{
		repo:   repo,
		logger: logger,
	}
}

func (svc *userService) FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error) {
	// 快路径
	obj, err := svc.repo.FindByWechatOpenId(ctx, info.OpenId)
	if err != repository.ErrUserNotFound {
		return obj, err
	}

	// 慢路径
	// 插入新用户
	u := domain.User{
		WeChat: info,
	}
	id, err := svc.repo.Create(ctx, u)
	if err != nil && err != repository.ErrUserDuplicate {
		return domain.User{}, err
	}

	// 这里还有一个问题：主从延迟，可能会查不到新插入的数据
	return svc.repo.FindById(ctx, id)
}

func (svc *userService) Login(ctx context.Context, obj domain.User) (domain.User, error) {
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

func (svc *userService) SignUp(ctx context.Context, obj domain.User) error {
	// 密码加密
	hash, err := bcrypt.GenerateFromPassword([]byte(obj.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	obj.Password = string(hash)

	_, err = svc.repo.Create(ctx, obj)
	return err
}

func (svc *userService) Profile(ctx context.Context, uid int64) (domain.User, error) {
	return domain.User{}, nil
}

func (svc *userService) FindOrCreateByPhone(ctx context.Context, phone string) (domain.User, error) {
	// 快路径
	obj, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		return obj, err
	}

	// 慢路径

	// 在系统资源不足，触发降级之后，不执行慢路径了
	//if ctx.Value("降级") == "true" {
	//	return domain.User{}, errors.New("系统降级了")
	//}

	// 插入新用户
	u := domain.User{
		Phone: phone,
	}
	id, err := svc.repo.Create(ctx, u)
	if err != nil && err != repository.ErrUserDuplicate {
		return domain.User{}, err
	}

	// 这里还有一个问题：主从延迟，可能会查不到新插入的数据
	return svc.repo.FindById(ctx, id)
}
