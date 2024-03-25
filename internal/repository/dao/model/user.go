package model

import "database/sql"

// User 持久化对象，对应数据库表结构
// 有些地方叫做 model 有些地方叫做 entity 还有地方叫做 po（ persistence object ）
type User struct {
	Id       int64          `gorm:"primaryKey;autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string
	Phone    sql.NullString `gorm:"unique"`

	WechatOpenId  sql.NullString
	WechatUnionId sql.NullString

	// 时间戳，单位毫秒
	CreateTime int64
	UpdateTime int64
}
