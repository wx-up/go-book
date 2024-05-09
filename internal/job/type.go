package job

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"github.com/wx-up/go-book/pkg/logger"
)

type Job interface {
	Name() string
	Run() error
}

// CronJobBuilder 用于给 Job 增强，增加了 logger 以及 prometheus 监控
type CronJobBuilder struct {
	logger logger.Logger
	vector *prometheus.SummaryVec
}

func NewCronJobBuilder(logger logger.Logger) *CronJobBuilder {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "learning",
		Subsystem: "go_book",
		Help:      "定时任务执行时间监控",
		Name:      "cron_job",
	}, []string{"job_name", "success"})
	prometheus.MustRegister(vector)
	return &CronJobBuilder{
		logger: logger,
		vector: vector,
	}
}

func (b *CronJobBuilder) Build(job Job) cron.Job {
	return cronJobFuncAdapter(func() error {
		startTime := time.Now()
		b.logger.Info("开始执行任务", logger.String("job_name", job.Name()))
		err := job.Run()
		duration := time.Since(startTime)
		b.logger.Info("任务执行完成", logger.String("job_name", job.Name()))
		if err != nil {
			b.logger.Error("任务运行失败", logger.String("job_name", job.Name()), logger.Error(err))
		}
		b.vector.WithLabelValues(job.Name(), strconv.FormatBool(err == nil)).Observe(float64(duration.Milliseconds()))
		return nil
	})
}

// cronJobFuncAdapter 适配器模式，用来适配 robfig 库的 cron.Job 接口
type cronJobFuncAdapter func() error

func (f cronJobFuncAdapter) Run() {
	_ = f()
}
