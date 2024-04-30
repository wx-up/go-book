package logger

import "go.uber.org/zap"

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(l *zap.Logger) *ZapLogger {
	return &ZapLogger{
		logger: l,
	}
}

func (z *ZapLogger) Debug(msg string, args ...Field) {
	z.logger.Debug(msg, z.toZapFields(args)...)
}

func (z *ZapLogger) Info(msg string, args ...Field) {
	// TODO implement me
	// panic("implement me")
}

func (z *ZapLogger) Warn(msg string, args ...Field) {
	// TODO implement me
	// panic("implement me")
}

func (z *ZapLogger) Error(msg string, args ...Field) {
	// TODO implement me
	// panic("implement me")
}

func (z *ZapLogger) Fatal(msg string, args ...Field) {
	// TODO implement me
	// panic("implement me")
}

func (z *ZapLogger) toZapFields(args []Field) []zap.Field {
	fs := make([]zap.Field, len(args))
	for i, arg := range args {
		fs[i] = zap.Any(arg.Key, arg.Value)
	}
	return fs
}
