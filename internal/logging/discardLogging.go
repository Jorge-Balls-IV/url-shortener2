package logging

import (
	"context"
	"log/slog"
)

type DiscardHandler struct{} //Хэндлекр, который ничего не логирует

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

func (dh *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}
func (dh *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}
func (dh *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return dh
}
func (dh *DiscardHandler) WithGroup(_ string) slog.Handler {
	return dh
}

// Логгер, который ничего не логирует
func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}
