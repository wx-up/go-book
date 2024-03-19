package cache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wx-up/go-book/internal/repository/cache/redismocks"

	"go.uber.org/mock/gomock"

	"github.com/redis/go-redis/v9"
)

func TestRedisCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name  string
		mock  func(ctrl *gomock.Controller) redis.Cmdable
		biz   string
		phone string
		code  string

		wantErr error
	}{
		{
			name: "验证码设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)

				// Eval 返回的是一个redis.Cmd对象，SetErr设置其错误信息，SetVal设置其返回值
				ret := &redis.Cmd{}
				ret.SetErr(nil)
				ret.SetVal(int64(0))
				cmd.EXPECT().
					Eval(gomock.Any(), luaSendCodeScript, []string{"code:login:13800138000"}, "123456").
					Return(ret)
				return cmd
			},
			biz:   "login",
			phone: "13800138000",
			code:  "123456",

			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cache := NewRedisCodeCache(tc.mock(ctrl))
			err := cache.Set(context.Background(), tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
