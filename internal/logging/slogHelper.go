package logging

import (
	"log/slog"
)
//переводим ошибку в атрибут slogger'a - набор ключ-значение
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
