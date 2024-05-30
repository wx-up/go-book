//go:build wireinject

package startup

import (
	"github.com/google/wire"
	grpcServer "github.com/wx-up/go-book/interactive/grpc"
	repository2 "github.com/wx-up/go-book/interactive/repository"
	cache2 "github.com/wx-up/go-book/interactive/repository/cache"
	dao2 "github.com/wx-up/go-book/interactive/repository/dao"
	"github.com/wx-up/go-book/interactive/service"
	service2 "github.com/wx-up/go-book/interactive/service"
)

var ThirdPartySet = wire.NewSet(
	InitTestRedis,
	InitTestMysql,
	CreateLogger,
)

var interactiveSvcSet = wire.NewSet(dao2.NewGORMInteractiveDAO,
	cache2.NewInteractiveRedisCache,
	repository2.NewCachedInteractiveRepository,
	service2.NewInteractiveService,
)

func InitInteractiveService() service.InteractiveService {
	wire.Build(ThirdPartySet, interactiveSvcSet)
	return service.NewInteractiveService(nil)
}

func InitInteractiveGRPCServer() *grpcServer.InteractiveServiceServer {
	wire.Build(ThirdPartySet, interactiveSvcSet, grpcServer.NewInteractiveServiceServer)
	return new(grpcServer.InteractiveServiceServer)
}
