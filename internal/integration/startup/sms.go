package startup

import (
	"github.com/wx-up/go-book/pkg/sms"
	"github.com/wx-up/go-book/pkg/sms/local"
)

// CreateLocalSMSService 集成测试就使用本地短信服务，不使用第三方，避免资产的消耗
func CreateLocalSMSService() sms.Service {
	return local.NewService()
}
