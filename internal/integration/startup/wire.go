//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/wx-up/go-book/internal/web"
	"github.com/wx-up/go-book/ioc"
)

func InitWebService() *gin.Engine {
	wire.Build(
		// 基础组件
		ThirdPartySet,

		UserHandlerSet,
		ArticleHandlerSet,
		WechatHandlerSet,
		JWTHandlerSet,

		// 中间件
		ioc.CreateMiddlewares,
		// web服务
		ioc.InitWeb,
	)

	return new(gin.Engine)
}

func CreateArticleHandler() *web.ArticleHandler {
	wire.Build(
		ThirdPartySet,
		ArticleHandlerSet,
	)
	return new(web.ArticleHandler)
}
