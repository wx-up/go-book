package ioc

import (
	"os"

	"github.com/wx-up/go-book/internal/service/oauth2/wechat"
)

func CreateOAuth2WechatService() wechat.Service {
	// 从环境变量中获取
	appId, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		panic("WECHAT_APP_ID not found in environment variables")
	}

	appSecret, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		panic("WECHAT_APP_SECRET not found in environment variables")
	}

	return wechat.NewService(appId, appSecret)
}
