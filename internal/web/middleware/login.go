package middleware

import (
	"net/http"
	"time"

	"github.com/wx-up/go-book/pkg/set"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// LoginMiddlewareBuilder builder 模式
type LoginMiddlewareBuilder struct {
	whiteList set.Set[string]
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{
		whiteList: set.NewMapSet[string](3),
	}
}

func (lm *LoginMiddlewareBuilder) IgnorePaths(paths ...string) *LoginMiddlewareBuilder {
	lm.whiteList.Add(paths...)
	return lm
}

func (lm *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if lm.whiteList.Exist(ctx.Request.URL.Path) {
			return
		}

		// 登陆验证
		sess := sessions.Default(ctx)
		if sess.Get("uid") == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		sess.Set("uid", sess.Get("uid"))

		// 刷新登陆状态
		updateTime := sess.Get("update_time")
		now := time.Now().UnixMilli()
		if updateTime == nil {
			sess.Options(sessions.Options{
				MaxAge: 30 * 60,
			})
			sess.Set("update_time", now)
			_ = sess.Save()
			return
		}

		// 每5秒刷新一次
		if now >= updateTime.(int64)+5*1000 {
			sess.Options(sessions.Options{
				MaxAge: 30 * 60,
			})
			sess.Set("update_time", now)
			_ = sess.Save()
			return
		}
	}
}
