package redisx

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/redis/go-redis/v9"
)

type PrometheusHook struct {
	vector *prometheus.SummaryVec
}

func NewPrometheusHook(opt prometheus.SummaryOpts) *PrometheusHook {
	// 动态标签由 PrometheusHook 自身决定
	// SummaryOpts 可以由外部传入
	vector := prometheus.NewSummaryVec(opt, []string{"name", "key_exist"})
	prometheus.MustRegister(vector)
	return &PrometheusHook{
		vector: vector,
	}
}

func (p *PrometheusHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

func (p *PrometheusHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		startTime := time.Now()
		var err error
		defer func() {
			val := ctx.Value("biz")
			_ = val
			duration := time.Since(startTime).Milliseconds()
			keyExist := !(err == redis.Nil)
			// cmd.Name() 执行的命令
			// 集群的情况下可以考虑使用 FullName
			p.vector.WithLabelValues(cmd.Name(), strconv.FormatBool(keyExist)).Observe(float64(duration))
		}()
		err = next(ctx, cmd)
		return err
	}
}

func (p *PrometheusHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	// TODO implement me
	panic("implement me")
}
