package dao

import (
	"context"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/wx-up/go-book/internal/repository/dao/model"
	"gorm.io/gorm"
)

var (
	ErrUserDuplicate = errors.New("邮箱冲突")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (obj model.User, err error) {
	// SELECT * FROM `users` WHERE `email` = ? LIMIT 1
	err = dao.db.WithContext(ctx).Where("email = ?", email).First(&obj).Error
	return
}

func (dao *UserDAO) FindByPhone(ctx context.Context, phone string) (obj model.User, err error) {
	// SELECT * FROM `users` WHERE `phone` = ? LIMIT 1
	err = dao.db.WithContext(ctx).Where("phone = ?", phone).First(&obj).Error
	return
}

func (dao *UserDAO) FindById(ctx context.Context, id int64) (obj model.User, err error) {
	// SELECT * FROM `users` WHERE `id` = ? LIMIT 1
	err = dao.db.WithContext(ctx).Where("id = ?", id).First(&obj).Error
	return
}

func (dao *UserDAO) Insert(ctx context.Context, obj model.User) (int64, error) {
	now := time.Now().UnixMilli()
	obj.CreateTime = now
	obj.UpdateTime = now
	err := dao.db.WithContext(ctx).Create(&obj).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		if mysqlErr.Number == 1062 {
			return 0, ErrUserDuplicate
		}
	}
	return obj.Id, err
}
