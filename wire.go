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

func InitWebService() *gin.Engine {
	wire.Build(
		// 基础组件
		ioc.CreateRedis,
		ioc.CreateMysql,

		// 用户服务
		service.NewUserService,
		repository.NewUserRepository,
		dao.NewUserDAO,
		cache.NewUserCache,

		// 短信服务
		code.NewSmsCodeService,
		ioc.CreateSMSService,
		repository.NewCodeRepository,
		cache.NewCodeCache,

		// user web
		web.NewUserHandler,

		// 中间件
		ioc.CreateMiddlewares,
		// web服务
		ioc.InitWeb,
	)

	return new(gin.Engine)
}
