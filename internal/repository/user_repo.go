package repository

import (
	"context"

	"github.com/wx-up/go-book/internal/repository/cache"

	"github.com/wx-up/go-book/internal/repository/dao/model"

	"github.com/wx-up/go-book/internal/repository/dao"

	"github.com/wx-up/go-book/internal/domain"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, cache *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: cache,
	}
}

// FindByEmail 用于登陆场景，根据 email 查找用户
// 登陆是一个比较低频的操作，有些网站对于登陆的 token 有很长的有效期，并且只要你一直活跃的话，token 时间还会不断刷新
// 因此没有必要设置缓存
func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := ur.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (ur *UserRepository) Create(ctx context.Context, u domain.User) error {
	return ur.dao.Insert(ctx, model.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (ur *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := ur.cache.Get(ctx, id)
	if err == nil {
		return u, nil
	}

	obj, err := ur.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	u = domain.User{
		Id:       obj.Id,
		Email:    obj.Email,
		Password: obj.Password,
	}

	// 写缓存失败的话，这里只是打日志，不处理错误
	// 缓存设置失败一般不是特别大的问题，可能是偶发性的，比如网络问题
	// 比如超时，因为用的同一个 ctx 如果数据库查询用了很多时间，那么缓存操作可能就超时了
	err = ur.cache.Set(ctx, u)
	if err != nil {
		// 打日志，做监控
	}
	return u, nil
}
