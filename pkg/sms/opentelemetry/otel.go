package opentelemetry

import (
	"context"

	"go.opentelemetry.io/otel/attribute"

	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/trace"

	"github.com/wx-up/go-book/pkg/sms"
)

type Service struct {
	svc    sms.Service
	tracer trace.Tracer
}

func NewService(svc sms.Service) *Service {
	prd := otel.GetTracerProvider()
	// tracer name 一般是 go mod 的名称+当前代码的相对路径
	// 比如当前项目的 go mod 名称为 github.com/wx-up/go-book
	// 当前代码的相对路径为 pkg/sms/opentelemetry
	// 所以整个 tracer name 为 github.com/wx-up/go-book/pkg/sms/opentelemetry
	tracer := prd.Tracer("github.com/wx-up/go-book/pkg/sms/opentelemetry")
	return &Service{
		svc:    svc,
		tracer: tracer,
	}
}

func (s *Service) Send(ctx context.Context, tplId string, params []sms.NameArg, phones ...string) error {
	ctx, span := s.tracer.Start(
		ctx,
		"send_sms",
		// 支持传递一些 option
		trace.WithSpanKind(trace.SpanKindClient),
	)

	// 注意 trace 没有 prometheus metric 笛卡尔积的问题
	span.SetAttributes(attribute.String("tpl", tplId))

	// span.End 也支持传递一些 option
	// defer span.End(trace.WithStackTrace(true))
	defer span.End()
	span.AddEvent("发送短信")
	err := s.svc.Send(ctx, tplId, params, phones...)
	if err != nil {
		span.RecordError(err)
	}
	return err
}

func (s *Service) Type() string {
	return "sms:otel"
}
