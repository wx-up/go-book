package domain

import (
	"time"

	"github.com/robfig/cron/v3"
)

type Job struct {
	Id       int64
	Name     string
	Cfg      string
	Executor string // 这个 job 的执行方式

	CancelFunc func() error

	Cron string // 这个 job 的 cron 表达式
}

var parser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

func (j *Job) NextTime() time.Time {
	// 在插入数据库的时候去检测 cron 表达式是否正确
	// 这里就直接忽略错误，认为一定是正确的
	s, _ := parser.Parse(j.Cron)
	return s.Next(time.Now())
}
