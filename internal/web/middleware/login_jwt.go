package middleware

import (
	"net/http"

	"github.com/wx-up/go-book/internal/web/jwt"

	"github.com/gin-gonic/gin"
)

// LoginJwtMiddlewareBuilder builder 模式
type LoginJwtMiddlewareBuilder struct {
	whiteList  []string
	jwtHandler jwt.Handler
}

func NewLoginJwtMiddlewareBuilder(handler jwt.Handler) *LoginJwtMiddlewareBuilder {
	return &LoginJwtMiddlewareBuilder{
		whiteList:  []string{"/users/login", "/users/signup"},
		jwtHandler: handler,
	}
}

func (lm *LoginJwtMiddlewareBuilder) IgnorePaths(paths ...string) *LoginJwtMiddlewareBuilder {
	lm.whiteList = append(lm.whiteList, paths...)
	return lm
}

func (lm *LoginJwtMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range lm.whiteList {
			if path == ctx.Request.URL.Path {
				// 直接 return 和 调用 ctx.Next() 之后再 return 效果一样的
				return
			}
		}
		var userClaim jwt.UserClaim
		err := lm.jwtHandler.ParseToken(ctx, &userClaim, jwt.AccessTokenKey)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if userClaim.Uid == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if userClaim.IsRefresh {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err = lm.jwtHandler.CheckSession(ctx, userClaim.Ssid)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 刷新 token
		// jwt 的 token 刷新即使生成了一个新的 token
		// 每隔10秒刷新一次
		//if userClaim.ExpiresAt.Sub(time.Now()) <= 50*time.Second {
		//	userClaim.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Second * 60))
		//	tokenStr, err := token.SignedString([]byte("go-book"))
		//	if err != nil {
		//		// 打印日志
		//		fmt.Println(err)
		//	}
		//	ctx.Header("x-jwt-token", tokenStr)
		//}

		ctx.Set("claims", userClaim)
	}
}
