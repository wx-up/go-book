package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wx-up/go-book/internal/domain"

	"github.com/wx-up/go-book/internal/web/jwt"

	svcmocks "github.com/wx-up/go-book/internal/service/mocks"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/require"

	"github.com/wx-up/go-book/internal/service"
	"go.uber.org/mock/gomock"
)

func TestArticleHandler_Publish(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) service.ArticleService
		req  PublishArticleReq

		wantCode int
		wantBody Result
	}{
		{
			name: "新建并发表，成功",
			req: PublishArticleReq{
				Title:   "test title",
				Content: "test content",
			},
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "test title",
					Content: "test content",
					Author: domain.Author{
						Id: 10,
					},
				}).Return(int64(1), nil)
				return svc
			},
			wantCode: http.StatusOK,
			wantBody: Result{
				Code: 0,
				Msg:  "发布成功",
				Data: map[string]any{
					"id": float64(1),
				},
			},
		},
		{
			name: "发表失败",
			req: PublishArticleReq{
				Title:   "test title",
				Content: "test content",
			},
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "test title",
					Content: "test content",
					Author: domain.Author{
						Id: 10,
					},
				}).Return(int64(0), errors.New("发表失败"))
				return svc
			},
			wantCode: http.StatusOK,
			wantBody: Result{
				Code: 5,
				Msg:  "服务器错误，请稍后再试",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			bs, err := json.Marshal(tc.req)
			require.NoError(t, err)
			req, err := http.NewRequest(
				http.MethodPost,
				"/articles/publish",
				bytes.NewReader(bs),
			)
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			gin.SetMode(gin.ReleaseMode)
			engine := gin.Default()
			engine.Use(func(context *gin.Context) {
				context.Set("claims", jwt.UserClaim{
					Uid: 10,
				})
			})
			ah := NewArticleHandler(tc.mock(ctrl), nil, nil)
			ah.RegisterRoutes(engine)
			engine.ServeHTTP(recorder, req)
			require.Equal(t, tc.wantCode, recorder.Code)
			if recorder.Code != http.StatusOK {
				return
			}
			var res Result
			err = json.NewDecoder(recorder.Body).Decode(&res)
			require.NoError(t, err)
			require.Equal(t, tc.wantBody, res)
		})
	}
}
