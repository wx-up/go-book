package dao

import "gorm.io/gorm"

type DBProvider func() *gorm.DB
