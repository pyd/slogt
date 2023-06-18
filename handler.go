package slogt

import (
	"context"

	"golang.org/x/exp/slog"
)

// a slog handler which captures logs for observation
// implements slog.Handler interface
type ObserverHandler struct {
	slog.Handler
}

func NewObserverHandler(handler slog.Handler) ObserverHandler {
	return ObserverHandler{
		handler,
	}
}

func (h ObserverHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h ObserverHandler) Handle(ctx context.Context, record slog.Record) error {
	return h.Handler.Handle(ctx, record)
}

func (h ObserverHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.Handler = h.Handler.WithAttrs(attrs)
	return h
}

func (h ObserverHandler) WithGroup(name string) slog.Handler {
	h.Handler = h.Handler.WithGroup(name)
	return h
}
