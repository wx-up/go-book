package tencent

import (
	"context"
	"fmt"

	typ "github.com/wx-up/go-book/pkg/sms"

	"github.com/hashicorp/go-multierror"

	"github.com/wx-up/go-book/pkg/slice"

	"github.com/wx-up/go-book/pkg"

	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	client   *sms.Client
	apdId    *string
	signName *string
}

// NewService 腾讯云的发送短信，需要一个 appId
func NewService(client *sms.Client, appId string, signName string) *Service {
	return &Service{
		client:   client,
		apdId:    pkg.ToPtr[string](appId),
		signName: pkg.ToPtr[string](signName),
	}
}

func (s *Service) Send(ctx context.Context, tplId string, params []typ.NameArg, phones ...string) error {
	resp, err := s.client.SendSmsWithContext(ctx, &sms.SendSmsRequest{
		SmsSdkAppId: s.apdId,
		SignName:    s.signName,
		TemplateParamSet: slice.Map[typ.NameArg, *string](params, func(idx int, val typ.NameArg) *string {
			return &val.Val
		}),
		TemplateId: pkg.ToPtr[string](tplId),
		PhoneNumberSet: slice.Map[string, *string](phones, func(idx int, val string) *string {
			return &val
		}),
	})
	if err != nil {
		return err
	}

	if resp.Response == nil {
		return fmt.Errorf("【腾讯云】短信发送失败，原因: 响应为空")
	}

	// 可以给多个手机号发送短信，每个手机号一个结果
	for _, status := range resp.Response.SendStatusSet {
		if status == nil {
			continue
		}
		if status.Code == nil {
			continue
		}
		if *(status.Code) != "Ok" {
			err = multierror.Append(err, fmt.Errorf("【腾讯云】短信发送失败，code: %s, 原因: %s", *status.Code, *status.Message))
		}
	}
	return err
}
