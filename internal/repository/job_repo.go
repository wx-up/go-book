package repository

import (
	"context"
	"time"

	"github.com/wx-up/go-book/internal/repository/dao/model"

	"github.com/wx-up/go-book/internal/repository/dao"

	"github.com/wx-up/go-book/internal/domain"
)

type JobRepository interface {
	Preempt(ctx context.Context) (domain.Job, error)
	Release(ctx context.Context, obj domain.Job) error
	UpdateTime(ctx context.Context, id int64) error
	UpdateNextTime(ctx context.Context, id int64, nextTime time.Time) error
	Stop(ctx context.Context, id int64) error
}

type CronJobRepository struct {
	dao dao.JobDAO
}

func (c *CronJobRepository) Stop(ctx context.Context, id int64) error {
	return c.dao.Stop(ctx, id)
}

func (c *CronJobRepository) UpdateNextTime(ctx context.Context, id int64, nextTime time.Time) error {
	return c.dao.UpdateNextTime(ctx, id, nextTime)
}

func (c *CronJobRepository) UpdateTime(ctx context.Context, id int64) error {
	return c.dao.UpdateTime(ctx, id)
}

func (c *CronJobRepository) Release(ctx context.Context, obj domain.Job) error {
	return c.dao.Release(ctx, model.Job{
		Id: obj.Id,
	})
}

func (c *CronJobRepository) Preempt(ctx context.Context) (domain.Job, error) {
	obj, err := c.dao.Preempt(ctx)
	if err != nil {
		return domain.Job{}, err
	}
	return c.toDomain(obj), nil
}

func (c *CronJobRepository) toDomain(obj model.Job) domain.Job {
	return domain.Job{
		Cfg: obj.Cfg,
	}
}
