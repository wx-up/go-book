package metric

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

type MiddlewareBuilder struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	// 有些企业是使用实例的IP地址，但是在 k8s 下Pod是会漂移的
	InstanceId string // 用于标记实例，可以从环境变量中获取
}

func (mb *MiddlewareBuilder) Build() gin.HandlerFunc {
	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      mb.Name + "_resp_time", // 使用下划线，中华线-会报错的
		Subsystem: mb.Subsystem,
		Namespace: mb.Namespace,
		Help:      mb.Help,
		ConstLabels: map[string]string{
			"instance_id": mb.InstanceId,
		},
		Objectives: map[float64]float64{
			0.5:  0.01,
			0.90: 0.01,
			0.95: 0.01,
			0.99: 0.001,
		},
	}, []string{"method", "route", "status"})
	prometheus.MustRegister(summaryVec)

	httpRequestCountGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      mb.Name + "_active_requests", // 使用下划线，中华线-会报错的
		Subsystem: mb.Subsystem,
		Namespace: mb.Namespace,
		Help:      mb.Help,
		ConstLabels: map[string]string{
			"instance_id": mb.InstanceId,
		},
	})
	prometheus.MustRegister(httpRequestCountGauge)

	return func(ctx *gin.Context) {
		start := time.Now()

		// ctx.Next() 执行后续业务逻辑，有可能 panic 因此这里使用 defer 保证一定会执行
		defer func() {
			httpRequestCountGauge.Inc()
			duration := time.Since(start)

			pattern := ctx.FullPath()
			// 路由没有匹配到，pattern 为空
			if pattern == "" {
				pattern = "unknown"
			}
			summaryVec.WithLabelValues(
				ctx.Request.Method,
				pattern,
				strconv.Itoa(ctx.Writer.Status()),
			).Observe(float64(duration.Milliseconds()))
		}()
		ctx.Next() // 执行后续的业务逻辑
	}
}
