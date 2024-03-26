package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/redis/go-redis/v9"

	"github.com/wx-up/go-book/internal/web"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
)

// LoginJwtMiddlewareBuilder builder 模式
type LoginJwtMiddlewareBuilder struct {
	whiteList []string
	cmd       redis.Cmdable
}

func NewLoginJwtMiddlewareBuilder(cmd redis.Cmdable) *LoginJwtMiddlewareBuilder {
	return &LoginJwtMiddlewareBuilder{
		whiteList: []string{"/users/login", "/users/signup"},
		cmd:       cmd,
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

		var userClaim web.UserClaim
		jwtStr := jwtSlice[1]

		token, err := jwt.ParseWithClaims(jwtStr, &userClaim, func(token *jwt.Token) (interface{}, error) {
			return []byte("go-book"), nil
		})
		// errors.Is(err, jwt.ErrTokenExpired)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if token == nil || !token.Valid || userClaim.Uid == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if userClaim.IsRefresh {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 查看 ssid 是否在redis中
		cnt, err := lm.cmd.Exists(ctx, fmt.Sprintf("users:sid:%s", userClaim.Ssid)).Result()
		if err != nil || cnt > 0 {
			// 要么 redis 出错，要么当前token已经退出登陆了

			// 这里其实可以考虑降级的
			// 如果 redis 出错了，不报错，直接继续执行，保证其他登陆用户不影响，毕竟退出登陆的用户在少数
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
