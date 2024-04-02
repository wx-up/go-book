package ioc

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/wx-up/go-book/pkg/ginx/middleware/logger"
	"go.uber.org/zap"

	"github.com/wx-up/go-book/internal/web/jwt"

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

func CreateJwtHandler(cmd redis.Cmdable) jwt.Handler {
	return jwt.NewRedisJwtHandler(cmd)
}

func CreateMiddlewares(jwtHandler jwt.Handler) []gin.HandlerFunc {
	accessLoggerBuilder := logger.NewBuilder(func(ctx *gin.Context, al *logger.AccessLog) {
		zap.L().Debug("HTTP请求", zap.Any("al", al))
	})
	// 这里监听配置变化
	viper.OnConfigChange(func(in fsnotify.Event) {
		allowReqBody := viper.GetBool("logger.allow_request_body")
		allowRespBody := viper.GetBool("logger.allow_response_body")
		accessLoggerBuilder.AllowRespBody(allowRespBody)
		accessLoggerBuilder.AllowReqBody(allowReqBody)
	})
	return []gin.HandlerFunc{
		// 请求体和响应打印
		accessLoggerBuilder.Build(),
		// 跨域
		cors.New(cors.Config{
			AllowMethods:     []string{"PUT", "PATCH", "POST"},
			AllowHeaders:     []string{"Authorization", "Content-Type"},
			ExposeHeaders:    []string{"X-Jwt-Token", "X-Refresh-Token"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				return true
			},
			MaxAge: 12 * time.Hour,
		}),

		// 登陆
		middleware.NewLoginJwtMiddlewareBuilder(jwtHandler).
			IgnorePaths("/users/code/send").
			IgnorePaths("/users/code/verify").
			IgnorePaths("/oauth2/wechat/callback").
			IgnorePaths("/users/refresh_token").
			IgnorePaths("/oauth2/wechat/auth_url").Build(),

		// 限流
		// ratelimit.NewBuilder(cmd, time.Second, 100).Build(),
	}
}
