package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// LoginMiddlewareBuilder builder 模式
type LoginMiddlewareBuilder struct {
	whiteList []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{
		whiteList: []string{"/users/login", "/users/signup"},
	}
}

func (lm *LoginMiddlewareBuilder) IgnorePaths(paths ...string) *LoginMiddlewareBuilder {
	lm.whiteList = append(lm.whiteList, paths...)
	return lm
}

func (lm *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range lm.whiteList {
			if path == ctx.Request.URL.Path {
				// 直接 return 和 调用 ctx.Next() 之后再 return 效果一样的
				return
			}
		}

		// 登陆验证
		sess := sessions.Default(ctx)
		if sess.Get("uid") == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
