package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
)

func Test_OpenTelemetry(t *testing.T) {
	res, err := newResource("demo", "v0.0.1")
	require.NoError(t, err)

	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// 初始化 trace provider
	// 这个 provider 就是用来在打点的时候构建 trace 的
	tp, err := newTraceProvider(res)
	require.NoError(t, err)
	// 用完需要关闭
	defer tp.Shutdown(context.Background())
	otel.SetTracerProvider(tp)

	server := gin.Default()
	server.Use(otelgin.Middleware("service"))
	server.GET("/test", func(ginCtx *gin.Context) {
		//// 名字唯一
		tracer := otel.Tracer("gitee.com/geekbang/basic-go/opentelemetry")
		var ctx context.Context = ginCtx
		ctx, span := tracer.Start(ctx, "top-span")
		defer span.End()
		time.Sleep(time.Second)
		span.AddEvent("发生了某事")
		ctx, subSpan := tracer.Start(ctx, "sub-span")
		defer subSpan.End()
		subSpan.SetAttributes(attribute.String("attr1", "value1"))
		time.Sleep(time.Millisecond * 300)
		ginCtx.String(http.StatusOK, "测试 span")
	})
	server.Run(":8082")
}
