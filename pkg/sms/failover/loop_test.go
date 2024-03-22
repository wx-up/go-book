package failover

import (
	"testing"
)

func TestLoopService_Send(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{
			name: "一次成功",
		},
		{
			name: "重试成功",
		},
		{
			name: "重试失败",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// TODO: 实现测试用例
		})
	}
}
