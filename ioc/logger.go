package ioc

import "go.uber.org/zap"

func CreateLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger
}
