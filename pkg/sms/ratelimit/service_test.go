package ratelimit

import (
	"testing"
)

func TestService_Send(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{
			name: "正常发送",
		},
		{
			name: "触发限流",
		},
		{
			name: "限流器异常",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// TODO: 实现测试用例
		})
	}
}
