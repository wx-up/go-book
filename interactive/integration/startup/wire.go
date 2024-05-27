//go:build wireinject

package startup

import (
	"github.com/google/wire"
	"github.com/wx-up/go-book/interactive/service"
)

var ThirdPartySet = wire.NewSet(
	InitTestRedis,
	InitTestMysql,
	CreateLogger,
)

func InitInteractiveService() service.InteractiveService {
	wire.Build(ThirdPartySet, interactiveSvcSet)
	return service.NewInteractiveService(nil)
}
