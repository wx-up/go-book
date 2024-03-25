package web

import (
	"net/http"

	"github.com/wx-up/go-book/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/wx-up/go-book/internal/service/oauth2/wechat"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	jwtHandler
}

var _ handler = (*OAuth2WechatHandler)(nil)

func NewOAuth2WechatHandler(svc wechat.Service, userSvc service.UserService) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:     svc,
		userSvc: userSvc,
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
	// 拿到临时授权码
	code := ctx.Query("code")
	// 拿到 state
	state := ctx.Query("state")

	// 获取 openId、unionId
	res, err := h.svc.Verify(ctx, code, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: -1,
			Msg:  "服务器错误",
			Data: nil,
		})
		return
	}

	// 查询用户信息
	u, err := h.userSvc.FindOrCreateByWechat(ctx, res)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: -1,
			Msg:  "服务器错误",
			Data: nil,
		})
		return
	}

	// 设置 jwt
	err = h.setJwtToken(ctx, u)
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
		Msg:  "登陆成功",
	})
}
