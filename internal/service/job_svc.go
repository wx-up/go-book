package service

import (
	"context"
	"time"

	"github.com/wx-up/go-book/pkg/logger"

	"github.com/wx-up/go-book/internal/repository"

	"github.com/wx-up/go-book/internal/domain"
)

// JobService 还需要暴露插入 job 的接口
// 正常流程：代码编写完成上线之后，调用插入 job 的接口新增一个 job 然后等待 job 的调度
type JobService interface {
	// Preempt 抢占
	Preempt(context.Context) (domain.Job, error)

	// Release 释放
	Release(context.Context, domain.Job) error
	ResetNextTime(context.Context, domain.Job) error

	Stop(context.Context, domain.Job) error
}

type CronJobService struct {
	repo            repository.JobRepository
	refreshInterval time.Duration
	logger          logger.Logger
}

func (m *CronJobService) Stop(ctx context.Context, job domain.Job) error {
	return m.repo.Stop(ctx, job.Id)
}

func (m *CronJobService) ResetNextTime(ctx context.Context, j domain.Job) error {
	nextTime := j.NextTime()

	// 任务没有下一次调度时间（ 一次性任务等等 ），则停止调度
	if nextTime.IsZero() || nextTime.Before(time.Now()) {
		return m.Stop(ctx, j)
	}
	return m.repo.UpdateNextTime(ctx, j.Id, nextTime)
}

func (m *CronJobService) Release(ctx context.Context, obj domain.Job) error {
	// TODO implement me
	panic("implement me")
}

func (m *CronJobService) Preempt(ctx context.Context) (domain.Job, error) {
	// 抢占任务
	j, err := m.repo.Preempt(ctx)
	if err != nil {
		return domain.Job{}, err
	}

	// 抢到任务之后，需要续约，保证任务不会被其他进程抢占
	ticker := time.NewTicker(m.refreshInterval)
	go func() {
		for range ticker.C {
			m.refresh(j)
		}
	}()

	// 抢占任务之后需要考虑释放的问题
	// 这里有两种方式实现，一种在接口中增加 Release(context.Context, int64) error 方法
	// 还有一种：在 job 中新增 cancelFunc 方法，用于释放，这里采用这个方法
	// 还有一种方式把 Preempt 改成这种签名 Preempt(context.Context) (domain.Job, func()error, error) 第二个参数就是释放函数
	j.CancelFunc = func() error {
		// 退出续约
		ticker.Stop()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// 这里可以考虑，如果下一次的执行时间和当前时间很近
		// 可以在内存中缓存这个 job，不进行 release 释放（ 减少数据库的压力 ），然后下次 preempt 时先从缓存中查找
		// 缓存的话，又要考虑缓存多少job，不能临近的job都缓存，可以考虑使用优先队列数据结构在内存中维持
		return m.repo.Release(ctx, j)
	}
	return j, nil
}

func (m *CronJobService) refresh(j domain.Job) {
	// 什么样的任务是续约失败的
	// 任务是 JobStatusPreempted 状态，但是 update_time 时间与当前时间的差值超过了 refresh_interval

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := m.repo.UpdateTime(ctx, j.Id)
	// 续约失败，其实做不了啥，因为你很难让当前任务的执行中断，只能打印日志
	if err != nil {
		// 可以考虑重试
		m.logger.Error("续约失败", logger.Error(err))
	}
}
