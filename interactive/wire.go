//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/wx-up/go-book/interactive/events/articles"
	grpcServer "github.com/wx-up/go-book/interactive/grpc"
	"github.com/wx-up/go-book/interactive/ioc"
	repository2 "github.com/wx-up/go-book/interactive/repository"
	cache2 "github.com/wx-up/go-book/interactive/repository/cache"
	dao2 "github.com/wx-up/go-book/interactive/repository/dao"
	"github.com/wx-up/go-book/interactive/service"
)

var ThirdPartySet = wire.NewSet(
	ioc.CreateRedis,
	ioc.CreateMysql,
	ioc.CreateLogger,
	ioc.InitKafka,
)

var InteractiveSvcSet = wire.NewSet(dao2.NewGORMInteractiveDAO,
	cache2.NewInteractiveRedisCache,
	repository2.NewCachedInteractiveRepository,
	service.NewInteractiveService,
)

func InitInteractiveGRPCServer() *grpcServer.InteractiveServiceServer {
	wire.Build(ThirdPartySet, InteractiveSvcSet, grpcServer.NewInteractiveServiceServer)
	return new(grpcServer.InteractiveServiceServer)
}

func InitApp() *App {
	wire.Build(
		ThirdPartySet,
		InteractiveSvcSet,
		grpcServer.NewInteractiveServiceServer,
		ioc.CreateGRPCServer,
		ioc.CreateConsumers,
		articles.NewReadEventKafkaConsumer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
