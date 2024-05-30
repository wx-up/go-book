package ioc

import (
	"github.com/wx-up/go-book/pkg/logger"
	"go.uber.org/zap"
)

func CreateLogger() logger.Logger {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(zapLogger)
}
