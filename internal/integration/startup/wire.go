//go:build wireinject

package startup

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
		InitTestRedis,
		InitTestMysql,

		// 用户服务
		service.NewUserService,
		repository.NewCacheUserRepository,
		dao.NewGORMUserDAO,
		cache.NewRedisUserCache,

		// 短信服务
		code.NewSmsCodeService,

		CreateLocalSMSService,

		repository.NewCacheCodeRepository,
		cache.NewRedisCodeCache,

		// user web
		web.NewUserHandler,

		// 中间件
		ioc.CreateMiddlewares,
		// web服务
		ioc.InitWeb,
	)

	return new(gin.Engine)
}
