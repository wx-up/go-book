package code

import (
	"context"
	"fmt"
	"math/rand"

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
	repo   *repository.CodeRepository
	tplId  string
}

func NewSmsCodeService(client sms.Service, repo *repository.CodeRepository, tplId string) *SmsCodeService {
	return &SmsCodeService{
		client: client,
		repo:   repo,
		tplId:  tplId,
	}
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
	err = s.client.Send(ctx, s.tplId, []sms.NameArg{{Name: "code", Val: code}}, phone)
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
