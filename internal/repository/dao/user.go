package dao

import (
	"context"
	"time"

	"github.com/wx-up/go-book/internal/repository/dao/model"
	"gorm.io/gorm"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(ctx context.Context, obj model.User) error {
	now := time.Now().UnixMilli()
	obj.CreateTime = now
	obj.UpdateTime = now
	return dao.db.WithContext(ctx).Create(&obj).Error
}
