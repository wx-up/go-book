package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wx-up/go-book/internal/service/oauth2/wechat"
)

type OAuth2WechatHandler struct {
	svc wechat.Service
}

var _ handler = (*OAuth2WechatHandler)(nil)

func NewOAuth2WechatHandler(svc wechat.Service) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc: svc,
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(engine *gin.Engine) {
	g := engine.Group("/oauth2/wechat")
	g.GET("/auth_url", h.AuthURL)
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	url, err := h.svc.AuthUrl(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: -1,
			Msg:  "服务器错误",
			Data: nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "success",
		Data: map[string]any{
			"auth_url": url,
		},
	})
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	// TODO: implement
}
