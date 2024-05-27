package startup

import (
	"github.com/google/wire"
	repository2 "github.com/wx-up/go-book/interactive/repository"
	cache2 "github.com/wx-up/go-book/interactive/repository/cache"
	dao2 "github.com/wx-up/go-book/interactive/repository/dao"
	service2 "github.com/wx-up/go-book/interactive/service"
)

var interactiveSvcSet = wire.NewSet(dao2.NewGORMInteractiveDAO,
	cache2.NewInteractiveRedisCache,
	repository2.NewCachedInteractiveRepository,
	service2.NewInteractiveService,
)
