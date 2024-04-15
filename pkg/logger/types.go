package logger

type Field struct {
	Key   string
	Value any
}

type NopLogger struct{}

func (n *NopLogger) Debug(msg string, args ...Field) {
	// TODO implement me
	panic("implement me")
}

func (n *NopLogger) Info(msg string, args ...Field) {
	// TODO implement me
	panic("implement me")
}

func (n *NopLogger) Warn(msg string, args ...Field) {
	// TODO implement me
	panic("implement me")
}

func (n *NopLogger) Error(msg string, args ...Field) {
	// TODO implement me
	panic("implement me")
}

func (n *NopLogger) Fatal(msg string, args ...Field) {
	// TODO implement me
	panic("implement me")
}

type Logger interface {
	Debug(msg string, args ...Field)
	Info(msg string, args ...Field)
	Warn(msg string, args ...Field)
	Error(msg string, args ...Field)
	Fatal(msg string, args ...Field)
}

func Error(val any) Field {
	return Field{
		Key:   "error",
		Value: val,
	}
}

func Int64(key string, val int64) Field {
	return Field{
		Key:   key,
		Value: val,
	}
}
