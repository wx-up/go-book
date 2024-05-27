// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
	"github.com/wx-up/go-book/internal/events/articles"
	"github.com/wx-up/go-book/internal/repository"
	"github.com/wx-up/go-book/internal/repository/cache"
	"github.com/wx-up/go-book/internal/repository/dao"
	"github.com/wx-up/go-book/internal/service"
	"github.com/wx-up/go-book/internal/service/code"
	"github.com/wx-up/go-book/internal/web"
	"github.com/wx-up/go-book/ioc"
	"github.com/wx-up/go-book/pkg/logger"
)

import (
	_ "github.com/spf13/viper/remote"
)

// Injectors from wire.go:

func InitWebService() *App {
	cmdable := ioc.CreateRedis()
	handler := ioc.CreateJwtHandler(cmdable)
	v := ioc.CreateMiddlewares(handler)
	db := ioc.CreateMysql()
	dbProvider := ioc.CreateDBProvider(db)
	userDAO := dao.NewGORMUserDAO(dbProvider)
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewCacheUserRepository(userDAO, userCache)
	zapLogger := ioc.CreateLogger()
	userService := service.NewUserService(userRepository, zapLogger)
	smsService := ioc.CreateSMSService()
	codeCache := cache.NewRedisCodeCache(cmdable)
	codeRepository := repository.NewCacheCodeRepository(codeCache)
	codeService := code.NewSmsCodeService(smsService, codeRepository)
	userHandler := web.NewUserHandler(userService, codeService, cmdable, handler)
	wechatService := ioc.CreateOAuth2WechatService()
	oAuth2WechatHandler := web.NewOAuth2WechatHandler(wechatService, userService, handler)
	gormArticleDAO := dao.NewGORMArticleDAO(dbProvider)
	cacheArticleRepository := repository.NewCacheArticleRepository(gormArticleDAO)
	articleService := service.NewArticleService(cacheArticleRepository)
	loggerZapLogger := logger.NewZapLogger(zapLogger)
	gormInteractiveDao := dao.NewGORMInteractiveDAO(db)
	redisInteractiveCache := cache.NewInteractiveRedisCache(cmdable)
	cacheInteractiveRepository := repository.NewCachedInteractiveRepository(gormInteractiveDao, redisInteractiveCache)
	interactiveService := service.NewInteractiveService(cacheInteractiveRepository)
	articleHandler := web.NewArticleHandler(articleService, interactiveService)
	engine := ioc.InitWeb(v, userHandler, oAuth2WechatHandler, articleHandler)
	client := ioc.InitKafka()
	readEventKafkaConsumer := article.NewReadEventKafkaConsumer(loggerZapLogger, cacheInteractiveRepository, client)
	v2 := ioc.CreateConsumers(readEventKafkaConsumer)
	app := &App{
		engine: engine,
		cs:     v2,
	}
	return app
}

// wire.go:

var thirdSet = wire.NewSet(ioc.CreateRedis, ioc.CreateMysql, ioc.CreateJwtHandler, ioc.CreateDBProvider, logger.NewZapLogger, ioc.CreateLogger, wire.Bind(new(logger.Logger), new(*logger.ZapLogger)))

var userSvcSet = wire.NewSet(service.NewUserService, repository.NewCacheUserRepository, dao.NewGORMUserDAO, cache.NewRedisUserCache)

var codeSvcSet = wire.NewSet(code.NewSmsCodeService, ioc.CreateSMSService, repository.NewCacheCodeRepository, cache.NewRedisCodeCache)

var userHandlerSet = wire.NewSet(web.NewUserHandler)

var wechatHandlerSet = wire.NewSet(web.NewOAuth2WechatHandler, ioc.CreateOAuth2WechatService)

var articleHandlerSet = wire.NewSet(web.NewArticleHandler, service.NewArticleService, service.NewInteractiveService, wire.Bind(new(repository.ArticleRepository), new(*repository.CacheArticleRepository)), repository.NewCacheArticleRepository, repository.NewCachedInteractiveRepository, wire.Bind(new(repository.InteractiveRepository), new(*repository.CachedInteractiveRepository)), dao.NewGORMInteractiveDAO, wire.Bind(new(dao.InteractiveDAO), new(*dao.GORMInteractiveDAO)), wire.Bind(new(dao.ArticleDAO), new(*dao.GORMArticleDAO)), dao.NewGORMArticleDAO, cache.NewInteractiveRedisCache, wire.Bind(new(cache.InteractiveCache), new(*cache.InteractiveRedisCache)))
