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
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
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
	err = dao.db.WithContext(ctx).Where("email = ?", email).First(&obj).Error
	return
}

func (dao *UserDAO) Insert(ctx context.Context, obj model.User) error {
	now := time.Now().UnixMilli()
	obj.CreateTime = now
	obj.UpdateTime = now
	err := dao.db.WithContext(ctx).Create(&obj).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		if mysqlErr.Number == 1062 {
			return ErrUserDuplicateEmail
		}
	}
	return err
}
