package local

import (
	"context"
	"fmt"

	typ "github.com/wx-up/go-book/pkg/sms"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Send(ctx context.Context, tplId string, params []typ.NameArg, phones ...string) error {
	for _, phone := range phones {
		fmt.Println("发送的手机号：", phone)
		fmt.Println(fmt.Sprintf("发送的短信内容：%+v", params))
	}
	return nil
}
