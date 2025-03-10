package zapx

import "go.uber.org/zap/zapcore"

type MyCore struct {
	zapcore.Core
}

func (c MyCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	for _, field := range fields {
		if field.Key == "phone" {
			phone := field.String
			field.String = phone[:3] + "****" + phone[7:]
		}
	}
	return c.Core.Write(entry, fields)
}
