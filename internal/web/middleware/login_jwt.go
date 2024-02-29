package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
)

// LoginJwtMiddlewareBuilder builder 模式
type LoginJwtMiddlewareBuilder struct {
	whiteList []string
}

func NewLoginJwtMiddlewareBuilder() *LoginJwtMiddlewareBuilder {
	return &LoginJwtMiddlewareBuilder{
		whiteList: []string{"/users/login", "/users/signup"},
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

		// 处理 jwt
		jwtOriStr := ctx.Request.Header.Get("Authorization")
		if jwtOriStr == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		jwtSlice := strings.SplitN(jwtOriStr, " ", 2)
		if len(jwtSlice) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		jwtStr := jwtSlice[1]
		token, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
			return []byte("go-book"), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if token == nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		println(token)
	}
}
