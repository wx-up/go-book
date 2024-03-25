package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/wx-up/go-book/internal/domain"
)

// Service 如果不需要单元测试的话，可以不用定义接口，直接定义结构体
type Service interface {
	AuthUrl(ctx context.Context, state string) (string, error)
	Verify(ctx context.Context, code string) (domain.WechatInfo, error)
}

type service struct {
	appId     string
	appSecret string
	client    *http.Client
}

func NewService(appId string, appSecret string) Service {
	return &service{
		appId:     appId,
		appSecret: appSecret,
		client:    http.DefaultClient,
	}
}

func NewServiceV1(appId string, appSecret string, client *http.Client) Service {
	return &service{
		appId:     appId,
		appSecret: appSecret,
		client:    client,
	}
}

func (s *service) AuthUrl(ctx context.Context, state string) (string, error) {
	const urlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	// haha.com 域名需要在微信开放平台中配置
	const redirectUri = "https://www.haha.com/oauth2/wechat/callback"
	return fmt.Sprintf(urlPattern, s.appId, url.QueryEscape(redirectUri), state), nil
}

func (s *service) Verify(ctx context.Context, code string) (domain.WechatInfo, error) {
	const requestUrl = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(requestUrl, s.appId, s.appSecret, code), nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}

	resp, err := s.client.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return domain.WechatInfo{}, err
	}

	// 读取响应
	var res Result
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		return domain.WechatInfo{}, fmt.Errorf("wechat auth failed, errcode: %d, errmsg: %s", res.ErrCode, res.ErrMsg)
	}

	return domain.WechatInfo{
		OpenId:  res.OpenId,
		UnionId: res.UnionId,
	}, nil
}

type Result struct {
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionId      string `json:"unionid"`
}
