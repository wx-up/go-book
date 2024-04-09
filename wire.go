//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/wx-up/go-book/internal/repository"
	"github.com/wx-up/go-book/internal/repository/cache"
	"github.com/wx-up/go-book/internal/repository/dao"
	"github.com/wx-up/go-book/internal/service"
	"github.com/wx-up/go-book/internal/service/code"
	"github.com/wx-up/go-book/internal/web"
	"github.com/wx-up/go-book/ioc"
)

var thirdSet = wire.NewSet(
	ioc.CreateRedis,
	ioc.CreateMysql,
	ioc.CreateLogger,
	ioc.CreateJwtHandler,
	ioc.CreateDBProvider,
)

var userSvcSet = wire.NewSet(
	service.NewUserService,
	repository.NewCacheUserRepository,
	dao.NewGORMUserDAO,
	cache.NewRedisUserCache,
)

var codeSvcSet = wire.NewSet(
	code.NewSmsCodeService,
	ioc.CreateSMSService,
	repository.NewCacheCodeRepository,
	cache.NewRedisCodeCache,
)

var userHandlerSet = wire.NewSet(
	web.NewUserHandler,
)

var wechatHandlerSet = wire.NewSet(
	web.NewOAuth2WechatHandler,
	ioc.CreateOAuth2WechatService,
)

var articleHandlerSet = wire.NewSet(
	web.NewArticleHandler,
	service.NewArticleService,
	wire.Bind(new(repository.ArticleRepository), new(*repository.CacheArticleRepository)),
	repository.NewCacheArticleRepository,
	wire.Bind(new(dao.ArticleDAO), new(*dao.GORMArticleDAO)),
	dao.NewGORMArticleDAO,
)

func InitWebService() *gin.Engine {
	wire.Build(
		// 基础组件
		thirdSet,
		// 用户服务
		userSvcSet,
		// 验证码服务
		codeSvcSet,

		// 用户 web
		userHandlerSet,
		// 微信 web
		wechatHandlerSet,
		// 文章 web
		articleHandlerSet,

		// 中间件
		ioc.CreateMiddlewares,
		// web服务
		ioc.InitWeb,
	)

	return new(gin.Engine)
}
