package ioc

import (
	"context"
	"time"

	"github.com/wx-up/go-book/internal/domain"
	"github.com/wx-up/go-book/internal/job"
	"github.com/wx-up/go-book/internal/service"
	"github.com/wx-up/go-book/pkg/logger"
)

func CreateJobScheduler(
	jobSvc service.JobService,
	logger logger.Logger,
	executors []job.Executor,
) *job.Scheduler {
	sch := job.NewScheduler(jobSvc, logger)
	for _, e := range executors {
		sch.RegisterExecutor(e)
	}
	return sch
}

func CreateJobExecutors(rankingSvc service.RankingService) []job.Executor {
	res := make([]job.Executor, 0, 1)
	funcExecutor := job.NewLocalFuncExecutor()
	funcExecutor.RegisterFunc("article_ranking", func(ctx context.Context, job domain.Job) error {
		ctx, cancel := context.WithTimeout(ctx, time.Second*30)
		defer cancel()
		return rankingSvc.TopN(ctx)
	})
	return append(res, funcExecutor)
}
