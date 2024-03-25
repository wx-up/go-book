//go:build manual

package wechat

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// manual 表示手动跑的测试，提前验证代码，主要验证与一些第三方的交互
func Test_service_Verify_manual(t *testing.T) {
	appId, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		panic("WECHAT_APP_ID not found in environment variables")
	}

	appSecret, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		panic("WECHAT_APP_SECRET not found in environment variables")
	}
	svc := NewService(appId, appSecret)
	res, err := svc.Verify(context.Background(), "code", "sate")
	require.NoError(t, err)
	t.Log(res)
}
