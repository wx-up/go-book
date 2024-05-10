package cache

import (
	"time"

	"github.com/ecodeclub/ekit"
)

type Cache interface {
	Set(context interface{}, key string, value any, expiration time.Duration) error
	Get(context interface{}, key string) (ekit.AnyValue, error)
}

type RedisCache struct{}

type LocalCache struct{}

type DoubleCache struct {
	local Cache
	redis Cache
}

func (d *DoubleCache) Set(context interface{}, key string, value any, expiration time.Duration) error {
	// TODO implement me
	panic("implement me")
}

func (d *DoubleCache) Get(context interface{}, key string) (ekit.AnyValue, error) {
	// TODO implement me
	panic("implement me")
}
