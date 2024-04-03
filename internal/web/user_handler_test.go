package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wx-up/go-book/internal/domain"

	"github.com/stretchr/testify/assert"

	svcmocks "github.com/wx-up/go-book/internal/service/mocks"

	"github.com/wx-up/go-book/internal/service"
	"github.com/wx-up/go-book/internal/service/code"
	"go.uber.org/mock/gomock"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/require"
)

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (service.UserService, code.Service)
		reqBody  func() []byte
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, code.Service) {
				userService := svcmocks.NewMockUserService(ctrl)
				userService.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "12@qq.com",
					Password: "ab12345!",
				}).Return(nil)
				codeService := svcmocks.NewMockService(ctrl)
				return userService, codeService
			},
			reqBody: func() []byte {
				bs, _ := json.Marshal(map[string]string{
					"password": "ab12345!",
					"email":    "12@qq.com",
				})
				return bs
			},
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			engine := gin.Default()
			mockUserService, mockCodeService := tt.mock(ctrl)
			userHandler := NewUserHandler(mockUserService, mockCodeService, nil, nil)
			userHandler.RegisterRoutes(engine)

			req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer(tt.reqBody()))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			engine.ServeHTTP(resp, req)

			assert.Equal(t, tt.wantCode, resp.Code)
			assert.Equal(t, tt.wantBody, resp.Body.String())
		})
	}
}
