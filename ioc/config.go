package ioc

import (
	"context"

	"github.com/spf13/viper"
)

type Config interface {
	GetString(ctx context.Context, key string) (string, error)
	MustGetString(ctx context.Context, key string) string
	GetStringOrDefault(ctx context.Context, key string, defaultValue string) string

	// 其他 GetXXX 接口

	// Unmarshal 接口
}

type ViperConfigAdapter struct {
	c *viper.Viper
}

func (v *ViperConfigAdapter) GetString(ctx context.Context, key string) (string, error) {
	// TODO implement me
	panic("implement me")
}

func (v *ViperConfigAdapter) MustGetString(ctx context.Context, key string) string {
	// TODO implement me
	panic("implement me")
}

func (v *ViperConfigAdapter) GetStringOrDefault(ctx context.Context, key string, defaultValue string) string {
	// TODO implement me
	panic("implement me")
}

type RedisConfig struct{}

func (f *RedisConfig) GetString(ctx context.Context, key string) (string, error) {
	// TODO implement me
	panic("implement me")
}

func (f *RedisConfig) MustGetString(ctx context.Context, key string) string {
	str, err := f.GetString(ctx, key)
	if err != nil {
		panic(err)
	}
	return str
}

func (f *RedisConfig) GetStringOrDefault(ctx context.Context, key string, defaultValue string) string {
	str, err := f.GetString(ctx, key)
	if err != nil {
		return defaultValue
	}
	return str
}
