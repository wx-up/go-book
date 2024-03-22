package auth

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"

	"github.com/wx-up/go-book/pkg/sms"
)

// Service 权限装饰器
type Service struct {
	sms sms.Service
	key string
}

func NewService(sms sms.Service) *Service {
	return &Service{
		sms: sms,
		key: "123456",
	}
}

type TokenClaim struct {
	jwt.RegisteredClaims
	TblId string
	Max   int64 // 最大发送数量
}

// Send 发送短信
// biz 为线下申请的静态token（ jwt token ） 包含了一些信息
func (s *Service) Send(ctx context.Context, biz string, params []sms.NameArg, phones ...string) error {
	// 解析 token
	var tc TokenClaim
	res, err := jwt.ParseWithClaims(biz, &tc, func(token *jwt.Token) (interface{}, error) {
		// 生成 jwtToken 时使用的 key
		return s.key, nil
	})
	if err != nil {
		return err
	}
	if res != nil && !res.Valid {
		return errors.New("token is invalid")
	}
	return s.sms.Send(ctx, tc.TblId, params, phones...)
}

func (s *Service) Type() string {
	// TODO implement me
	panic("implement me")
}
