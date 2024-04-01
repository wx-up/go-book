package zapx

import (
	"go.uber.org/zap/zapcore"
)

type SensitiveCore struct {
	zapcore.Core
}

func (c *SensitiveCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	for i, fd := range fields {
		if fd.Key == "password" {
			fd.String = "*****"
		}
		if fd.Key == "phone" {
			fd.String = fd.String[:3] + "****" + fd.String[7:]
		}
		fields[i] = fd
	}
	return c.Core.Write(ent, fields)
}
