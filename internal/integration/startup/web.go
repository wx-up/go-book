package startup

import (
	"github.com/google/wire"
	"github.com/wx-up/go-book/internal/repository"
	"github.com/wx-up/go-book/internal/repository/cache"
	"github.com/wx-up/go-book/internal/repository/dao"
	"github.com/wx-up/go-book/internal/service"
	"github.com/wx-up/go-book/internal/service/code"
	"github.com/wx-up/go-book/internal/web"
	"github.com/wx-up/go-book/internal/web/jwt"
	"gorm.io/gorm"
)

var ThirdPartySet = wire.NewSet(
	InitTestRedis,
	InitTestMysql,
	CreateLogger,
)

var ArticleHandlerSet = wire.NewSet(
	web.NewArticleHandler,
	wire.Bind(new(repository.ArticleRepository), new(*repository.CacheArticleRepository)),
	service.NewArticleService,
	repository.NewCacheArticleRepository,
	CreateArticleDAO,
)

func CreateArticleDAO(db *gorm.DB) dao.ArticleDAO {
	return dao.NewGORMArticleDAO(func() *gorm.DB {
		return db
	})
}

var UserHandlerSet = wire.NewSet(
	web.NewUserHandler,
	service.NewUserService,
	code.NewSmsCodeService,
	CreateLocalSMSService,
	repository.NewCacheUserRepository,
	cache.NewRedisUserCache,
	CreateUserDAO,
	repository.NewCacheCodeRepository,
	cache.NewRedisCodeCache,
)

func CreateUserDAO(db *gorm.DB) dao.UserDAO {
	return dao.NewGORMUserDAO(func() *gorm.DB {
		return db
	})
}

var WechatHandlerSet = wire.NewSet(
	web.NewOAuth2WechatHandler,
	CreateOAuth2WechatService,
)

var JWTHandlerSet = wire.NewSet(
	wire.Bind(new(jwt.Handler), new(*jwt.RedisJwtHandler)),
	jwt.NewRedisJwtHandler,
)

var interactiveSvcSet = wire.NewSet(dao.NewGORMInteractiveDAO,
	cache.NewInteractiveRedisCache,
	repository.NewCachedInteractiveRepository,
	service.NewInteractiveService,
)
