package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/wx-up/go-book/internal/integration/startup"

	"github.com/wx-up/go-book/internal/web"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUserHandler_SendCode(t *testing.T) {
	server := startup.InitWebService()
	redisClient := startup.InitTestRedis()
	testCases := []struct {
		name string

		phone string

		before func(t *testing.T) // 准备数据
		after  func(t *testing.T) // 数据验证以及数据清理（ 删除当前测试产生的数据，避免影响其他测试用例 ）

		wantCode int
		wantRes  web.Result
	}{
		{
			name:  "发送成功",
			phone: "15658283276",
			before: func(t *testing.T) {
				// 不需要准备数据
			},
			after: func(t *testing.T) {
				key := fmt.Sprintf("code:%s:%s", "login", "15658283276")
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				v, err := redisClient.Get(ctx, key).Result()
				require.NoError(t, err)
				// 因为验证码的生成逻辑是内部写死的，因此这里只能比较长度，不能比较实际的值
				assert.Equal(t, 6, len(v))

				// 验证有效期
				dur, err := redisClient.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, dur > time.Minute*9+time.Second*50)

				// 删除测试数据，保证其他测试用例不受影响
				err = redisClient.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			wantCode: http.StatusOK,
			wantRes: web.Result{
				Code: 0,
				Msg:  "验证码发送成功",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			req, err := http.NewRequest(http.MethodPost,
				"/users/code/send",
				bytes.NewBuffer([]byte(fmt.Sprintf(fmt.Sprintf(`{"phone":"%s"}`, tc.phone)))),
			)
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()

			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)

			var res web.Result
			// json.Unmarshal(resp.Body.Bytes(), &res)
			// 下面的方式性能好一点，上面的方式 resp.Body.Bytes() 读取了一次内容 json.Unmarshal 又读取一次，相当于读取了两次
			// 下面的话只会被 json.NewDecoder 读取一次
			err = json.NewDecoder(resp.Body).Decode(&res)
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)
			tc.after(t)
		})
	}
}
