package model

type AsyncSms struct {
	Id int64 `gorm:"column:id;primary_key"`

	// 重试次数
	RetryCnt int
	// 最大重试次数
	RetryMax int

	Status     int8
	CreateTime int64
	UpdateTime int64 `gorm:"column:update_time;index"`
}
