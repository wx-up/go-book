package domain

import "time"

// User 领域对象
// 在 DDD 中也叫做 entity 实体
// 其他也有地方叫做 BO（ business object ）
type User struct {
	Id         int64
	Nickname   string
	Email      string
	Password   string
	Phone      string
	CreateTime time.Time

	WeChat WechatInfo
}
