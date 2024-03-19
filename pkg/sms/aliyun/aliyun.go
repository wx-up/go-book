package aliyun

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wx-up/go-book/pkg"

	sms "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	typ "github.com/wx-up/go-book/pkg/sms"
)

type Service struct {
	client   *sms.Client
	signName string
}

func (s *Service) Type() string {
	return "aliyun"
}

func NewService(client *sms.Client, signName string) *Service {
	return &Service{
		client:   client,
		signName: signName,
	}
}

func (s *Service) Send(ctx context.Context, tplId string, params []typ.NameArg, phones ...string) error {
	req := &sms.SendSmsRequest{
		// 阿里云多个手机号用逗号分割
		PhoneNumbers: pkg.ToPtr[string](strings.Join(phones, ",")),
		SignName:     pkg.ToPtr[string](s.signName),
		TemplateCode: pkg.ToPtr[string](tplId),
	}
	argsMap := make(map[string]string, len(params))
	for _, arg := range params {
		argsMap[arg.Name] = arg.Val
	}
	tplParamBytes, err := json.Marshal(argsMap)
	if err != nil {
		return err
	}
	req.TemplateParam = pkg.ToPtr[string](string(tplParamBytes))
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	if resp.Body == nil {
		return fmt.Errorf("【阿里云】短信发送失败，原因: 响应为空")
	}
	if *resp.Body.Code != "OK" {
		return fmt.Errorf("【阿里云】短信发送失败，原因: %s", *resp.Body.Message)
	}
	return nil
}
