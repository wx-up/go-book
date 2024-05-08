package otelx

import (
	"context"
	"sync/atomic"

	"go.opentelemetry.io/otel/trace/noop"

	"go.opentelemetry.io/otel/trace/embedded"

	"go.opentelemetry.io/otel/trace"
)

type MyTraceProvider struct {
	embedded.TracerProvider
	Enabled     *atomic.Bool
	provider    trace.TracerProvider
	nopProvider trace.TracerProvider
}

func NewMyTraceProvider(tp trace.TracerProvider) *MyTraceProvider {
	return &MyTraceProvider{
		Enabled:     &atomic.Bool{},
		provider:    tp,
		nopProvider: noop.NewTracerProvider(),
	}
}

func (m *MyTraceProvider) Tracer(name string, options ...trace.TracerOption) trace.Tracer {
	if m.Enabled.Load() {
		return m.provider.Tracer(name, options...)
	}
	return m.nopProvider.Tracer(name, options...)
}

func (m *MyTraceProvider) TracerV1(name string, options ...trace.TracerOption) *MyTracer {
	return &MyTracer{
		tracer:    m.provider.Tracer(name, options...),
		nopTracer: m.nopProvider.Tracer(name, options...),
		Enabled:   m.Enabled,
	}
}

type MyTracer struct {
	embedded.Tracer
	Enabled   *atomic.Bool
	tracer    trace.Tracer
	nopTracer trace.Tracer
}

func (m *MyTracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	if m.Enabled.Load() {
		return m.tracer.Start(ctx, spanName, opts...)
	}
	return m.nopTracer.Start(ctx, spanName, opts...)
}
