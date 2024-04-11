package ginx

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"

	"github.com/wx-up/go-book/pkg/logger"

	"github.com/gin-gonic/gin"
)

var L logger.Logger = &logger.NopLogger{}

func InitLogger(l logger.Logger) {
	L = l
}

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func WrapHandleWithClaim[Req any, Claim jwt.Claims](
	claimKey string,
	handle func(ctx *gin.Context, req Req, claim Claim) (Result, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.ShouldBind(&req); err != nil {
			ctx.JSON(http.StatusOK, Result{
				Code: -1,
				Msg:  "参数错误",
			})
			return
		}
		val, ok := ctx.Get(claimKey)
		if !ok {
			ctx.JSON(http.StatusOK, Result{
				Code: -1,
				Msg:  "用户未登录",
			})
			return
		}
		claim, ok := val.(Claim)
		if !ok {
			ctx.JSON(http.StatusOK, Result{
				Code: -1,
				Msg:  "用户未登录",
			})
			return
		}
		result, err := handle(ctx, req, claim)
		if err != nil {
			// 统一打日志
			L.Error("业务逻辑处理错误", logger.Field{
				Key:   "err",
				Value: err,
			})
		}
		ctx.JSON(http.StatusOK, result)
	}
}

func WrapHandleV2[Req any](
	f func(ctx *gin.Context, req Req) (Result, error),
	before func(),
	after func(),
) gin.HandlerFunc {
	return func(context *gin.Context) {
		var req Req
		before()
		_, _ = f(context, req)
		after()
	}
}

func WrapHandle[Req any](f func(ctx *gin.Context, req Req) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.ShouldBind(&req); err != nil {
			ctx.JSON(http.StatusOK, Result{
				Code: -1,
				Msg:  "参数错误",
			})
			return
		}
		result, err := f(ctx, req)
		if err != nil {
			L.Error("业务逻辑处理错误", logger.Field{
				Key:   "err",
				Value: err,
			})
		}
		ctx.JSON(http.StatusOK, result)
	}
}
