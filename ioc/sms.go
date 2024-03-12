package ioc

import (
	"github.com/wx-up/go-book/pkg/sms"
	"github.com/wx-up/go-book/pkg/sms/local"
)

func CreateSMSService() sms.Service {
	// 本地短信服务
	return local.NewService()
}
