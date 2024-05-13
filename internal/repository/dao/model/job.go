package model

type Job struct {
	Id  int64
	Cfg string

	Status int8

	Version int64

	// 下一次被调度的时间
	// 这里和 status 建立一个联合索引效果会更好，因此 preempt 的时候根据 status 和 next_time 过滤
	NextTime int64 `gorm:"index"`

	UpdateTime int64
	CreateTime int64
}

const (
	JobStatusWaiting   = 0
	JobStatusPreempted = 1 // 已经被抢占
	JobStatusPaused    = 2 // 暂停调度
)
