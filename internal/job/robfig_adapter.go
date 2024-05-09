package job

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/wx-up/go-book/pkg/logger"
)

type RobfigJobAdapter struct {
	job Job
	l   logger.Logger
	p   *prometheus.SummaryVec
}

func NewRobfigJobAdapter(job Job, l logger.Logger) *RobfigJobAdapter {
	sum := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		ConstLabels: map[string]string{
			"job_name": job.Name(),
		},
	}, []string{"success"})
	prometheus.MustRegister(sum)
	return &RobfigJobAdapter{
		job: job,
		l:   l,
		p:   sum,
	}
}

func (r *RobfigJobAdapter) Run() {
	startTime := time.Now()
	err := r.job.Run()
	duration := time.Since(startTime)
	if err != nil {
		r.l.Error("任务运行失败", logger.Error(err), logger.String("job", r.job.Name()))
	}
	r.p.WithLabelValues(strconv.FormatBool(err == nil)).Observe(float64(duration.Milliseconds()))
}
