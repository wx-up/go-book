package ioc

import (
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wx-up/go-book/internal/web/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wx-up/go-book/internal/web"
)

func InitWeb(ms []gin.HandlerFunc,
	uh *web.UserHandler,
	wh *web.OAuth2WechatHandler,
) *gin.Engine {
	engine := gin.Default()
	engine.Use(ms...)
	uh.RegisterRoutes(engine)
	wh.RegisterRoutes(engine)
	return engine
}

func CreateMiddlewares(cmd redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		// 跨域
		cors.New(cors.Config{
			AllowMethods:     []string{"PUT", "PATCH", "POST"},
			AllowHeaders:     []string{"Authorization", "Content-Type"},
			ExposeHeaders:    []string{"X-Jwt-Token"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				return true
			},
			MaxAge: 12 * time.Hour,
		}),

		// 登陆
		middleware.NewLoginJwtMiddlewareBuilder().
			IgnorePaths("/users/code/send").
			IgnorePaths("/users/code/verify").
			IgnorePaths("/oauth2/wechat/callback").
			IgnorePaths("/oauth2/wechat/auth_url").Build(),

		// 限流
		// ratelimit.NewBuilder(cmd, time.Second, 100).Build(),
	}
}
