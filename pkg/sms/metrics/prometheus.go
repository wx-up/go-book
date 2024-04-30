package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/wx-up/go-book/pkg/sms"
)

type PrometheusDecorator struct {
	svc    sms.Service
	vector *prometheus.SummaryVec
}

func NewPrometheusDecorator(svc sms.Service) sms.Service {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  "wx",
		Subsystem:  "go_book",
		Name:       "sms_resp_time",
		Help:       "统计 SMS 服务的性能数据",
		Objectives: map[float64]float64{},
		// 动态标签，这里千万不要把 phone 也加入进来，还是那个问题，这会导致时许数据很多，到时候就炸裂了
	}, []string{"name"})
	prometheus.MustRegister(vector)
	return &PrometheusDecorator{
		svc:    svc,
		vector: vector,
	}
}

func (p *PrometheusDecorator) Send(ctx context.Context, tplId string, params []sms.NameArg, phones ...string) error {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime).Milliseconds()
		p.vector.WithLabelValues(p.svc.Type()).Observe(float64(duration))
	}()
	return p.svc.Send(ctx, tplId, params, phones...)
}

func (p *PrometheusDecorator) Type() string {
	return "sms:metric:prometheus"
}
