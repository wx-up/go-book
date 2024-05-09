package job

import (
	"github.com/robfig/cron/v3"
	"github.com/wx-up/go-book/pkg/logger"
)

func InitJobs(l logger.Logger) *cron.Cron {
	c := cron.New(cron.WithSeconds())
	cronBuilder := NewCronJobBuilder(l)
	_, err := c.AddJob("@every 1m", cronBuilder.Build(NewArticleRankingJob()))
	if err != nil {
		panic(err)
	}
	return c
}
