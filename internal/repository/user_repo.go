package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/wx-up/go-book/internal/repository/cache"

	"github.com/wx-up/go-book/internal/repository/dao/model"

	"github.com/wx-up/go-book/internal/repository/dao"

	"github.com/wx-up/go-book/internal/domain"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	Create(ctx context.Context, u domain.User) (int64, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
}

type CacheUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewCacheUserRepository(dao dao.UserDAO, cache cache.UserCache) UserRepository {
	return &CacheUserRepository{
		dao:   dao,
		cache: cache,
	}
}

// FindByEmail 用于登陆场景，根据 email 查找用户
// 登陆是一个比较低频的操作，有些网站对于登陆的 token 有很长的有效期，并且只要你一直活跃的话，token 时间还会不断刷新
// 因此没有必要设置缓存
func (ur *CacheUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := ur.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return ur.modelToDomain(u), nil
}

func (ur *CacheUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := ur.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return ur.modelToDomain(u), nil
}

func (ur *CacheUserRepository) Create(ctx context.Context, u domain.User) (int64, error) {
	return ur.dao.Insert(ctx, ur.domainToModel(u))
}

func (ur *CacheUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := ur.cache.Get(ctx, id)
	if err == nil {
		return u, nil
	}

	obj, err := ur.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	u = ur.modelToDomain(obj)

	// 写缓存失败的话，这里只是打日志，不处理错误
	// 缓存设置失败一般不是特别大的问题，可能是偶发性的，比如网络问题
	// 比如超时，因为用的同一个 ctx 如果数据库查询用了很多时间，那么缓存操作可能就超时了
	err = ur.cache.Set(ctx, u)
	if err != nil {
		// 打日志，做监控
	}
	return u, nil
}

func (ur *CacheUserRepository) domainToModel(u domain.User) model.User {
	obj := model.User{
		Id:       u.Id,
		Email:    sql.NullString{String: u.Email, Valid: u.Email != ""},
		Password: u.Password,
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
	}
	if !u.CreateTime.IsZero() {
		obj.CreateTime = u.CreateTime.UnixMilli()
	}
	return obj
}

func (ur *CacheUserRepository) modelToDomain(u model.User) domain.User {
	return domain.User{
		Id:         u.Id,
		Email:      u.Email.String,
		Password:   u.Password,
		Phone:      u.Phone.String,
		CreateTime: time.UnixMilli(u.CreateTime),
	}
}
