package domain

// User 领域对象
// 在 DDD 中也叫做 entity 实体
// 其他也有地方叫做 BO（ business object ）
type User struct {
	Email    string
	Password string
}
