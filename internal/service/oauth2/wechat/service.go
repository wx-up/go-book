package wechat

import (
	"context"
	"fmt"
	"net/url"

	"github.com/lithammer/shortuuid/v4"
)

// Service 如果不需要单元测试的话，可以不用定义接口，直接定义结构体
type Service interface {
	AuthUrl(ctx context.Context) (string, error)
}

type service struct {
	appId string
}

func NewService(appId string) Service {
	return &service{
		appId: appId,
	}
}

func (s *service) AuthUrl(ctx context.Context) (string, error) {
	const urlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	// haha.com 域名需要在微信开放平台中配置
	const redirectUri = "https://haha.com/oauth2/wechat/callback"

	// shortuuid.New() 生成的 uuid 就比较短，冲突的概率会比较高
	return fmt.Sprintf(urlPattern, s.appId, url.QueryEscape(redirectUri), shortuuid.New()), nil
}
