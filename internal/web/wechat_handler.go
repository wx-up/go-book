package web

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/wx-up/go-book/internal/domain"
	ijwt "github.com/wx-up/go-book/internal/web/jwt"

	"github.com/lithammer/shortuuid/v4"

	"github.com/wx-up/go-book/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/wx-up/go-book/internal/service/oauth2/wechat"
)

type OAuth2WechatHandler struct {
	svc        wechat.Service
	userSvc    service.UserService
	jwtHandler ijwt.Handler
}

var _ handler = (*OAuth2WechatHandler)(nil)

func NewOAuth2WechatHandler(svc wechat.Service, userSvc service.UserService, jwtHandler ijwt.Handler) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:        svc,
		userSvc:    userSvc,
		jwtHandler: jwtHandler,
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(engine *gin.Engine) {
	g := engine.Group("/oauth2/wechat")
	g.GET("/auth_url", h.AuthURL)
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	state := shortuuid.New()
	url, err := h.svc.AuthUrl(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: -1,
			Msg:  "服务器错误",
			Data: nil,
		})
		return
	}
	err = h.setStateCookie(ctx, state)
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

func (h *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, StateClaim{
		State: state,
		// 整个扫码流程 1 分钟足够
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	})
	t, err := token.SignedString([]byte("state"))
	if err != nil {
		return fmt.Errorf("生成 state jwt 失败：%w", err)
	}
	ctx.SetCookie("jwt-state", t, 60, "/oauth2/wechat/callback", "", false, true)
	return nil
}

type StateClaim struct {
	jwt.RegisteredClaims
	State string
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	// 拿到临时授权码
	code := ctx.Query("code")
	err := h.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: -1,
			Msg:  "登陆失败",
			Data: nil,
		})
		return
	}

	// 获取 openId、unionId
	res, err := h.svc.Verify(ctx, code)
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
	err = h.setLoginToken(ctx, u)
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

func (h *OAuth2WechatHandler) setLoginToken(ctx *gin.Context, u domain.User) error {
	err := h.jwtHandler.SetAccessToken(ctx, u.Id)
	if err != nil {
		return err
	}
	err = h.jwtHandler.SetRefreshToken(ctx, u.Id)
	if err != nil {
		return err
	}
	return nil
}

func (h *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	// 拿到 state
	state := ctx.Query("state")
	// 拿到 cookie
	t, err := ctx.Cookie("jwt-state")
	if err != nil {
		// 记录日志做好监控
		ctx.JSON(http.StatusOK, Result{
			Code: -1,
			Msg:  "登陆失败",
			Data: nil,
		})
		return fmt.Errorf("获取 state cookie 错误：%w", err)
	}
	var c StateClaim
	tRes, err := jwt.ParseWithClaims(t, &c, func(token *jwt.Token) (interface{}, error) {
		return []byte("state"), nil
	})
	if err != nil || !tRes.Valid {
		return fmt.Errorf("state 无效")
	}
	if c.State != state {
		return errors.New("state 不相等")
	}
	return nil
}
