// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	"github.com/google/wire"
	"github.com/wx-up/go-book/interactive/grpc"
	"github.com/wx-up/go-book/interactive/repository"
	"github.com/wx-up/go-book/interactive/repository/cache"
	"github.com/wx-up/go-book/interactive/repository/dao"
	"github.com/wx-up/go-book/interactive/service"
)

// Injectors from wire.go:

func InitInteractiveService() service.InteractiveService {
	db := InitTestMysql()
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	cmdable := InitTestRedis()
	interactiveCache := cache.NewInteractiveRedisCache(cmdable)
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, interactiveCache)
	interactiveService := service.NewInteractiveService(interactiveRepository)
	return interactiveService
}

func InitInteractiveGRPCServer() *grpc.InteractiveServiceServer {
	db := InitTestMysql()
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	cmdable := InitTestRedis()
	interactiveCache := cache.NewInteractiveRedisCache(cmdable)
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, interactiveCache)
	interactiveService := service.NewInteractiveService(interactiveRepository)
	interactiveServiceServer := grpc.NewInteractiveServiceServer(interactiveService)
	return interactiveServiceServer
}

// wire.go:

var ThirdPartySet = wire.NewSet(
	InitTestRedis,
	InitTestMysql,
	CreateLogger,
)

var interactiveSvcSet = wire.NewSet(dao.NewGORMInteractiveDAO, cache.NewInteractiveRedisCache, repository.NewCachedInteractiveRepository, service.NewInteractiveService)
