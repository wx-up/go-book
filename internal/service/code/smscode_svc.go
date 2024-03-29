package code

import (
	"context"
	"fmt"
	"math/rand"

	"go.uber.org/atomic"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/wx-up/go-book/internal/repository"

	"github.com/wx-up/go-book/pkg/sms"
)

var (
	ErrCodeSendTooMany   = repository.ErrCodeSendTooMany
	ErrCodeVerifyTooMany = repository.ErrCodeVerifyTooMany
)

// SmsCodeService 基于短信的验证码服务
type SmsCodeService struct {
	client sms.Service
	repo   repository.CodeRepository
	tplId  atomic.String
}

func NewSmsCodeService(client sms.Service, repo repository.CodeRepository) Service {
	svc := &SmsCodeService{
		client: client,
		repo:   repo,
	}
	svc.tplId.Store("123")
	viper.OnConfigChange(func(in fsnotify.Event) {
		// 配置信息发送变化
		svc.tplId.Store(viper.GetString("sms.tpl_id"))
	})
	return svc
}

func (s *SmsCodeService) Send(ctx context.Context, biz string, phone string) error {
	// 生成验证码
	code := s.generateCode()
	// 保存验证码
	err := s.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	// 发送验证码
	err = s.client.Send(ctx, s.tplId.Load(), []sms.NameArg{{Name: "code", Val: code}}, phone)
	if err != nil {
		// 这里不应该删除 redis 的 key
		// 因为错误有可能是超时等问题，你不知道验证码是否真的发送
	}
	return err
}

func (s *SmsCodeService) generateCode() string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}

func (s *SmsCodeService) Verify(ctx context.Context, biz string, phone string, code string) error {
	return s.repo.Verify(ctx, biz, phone, code)
}
