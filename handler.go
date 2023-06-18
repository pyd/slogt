package slogt

import (
	"context"
	"errors"
	"io"

	"golang.org/x/exp/slog"
)

// ObserverHandler accepts an observer only if it implements this interface
// see ObserverHandler.Handle()
type RecordCollector interface {
	addRecord(slog.Record)
}

// a handler - implementing the slog.Handler interface - for the slog logger
// additionally it provides logs received from the slog Logger to the observer
type ObserverHandler struct {
	slog.Handler
	RecordCollector
}

// ObserverHandler constructor
// an error is returned if handler or observer arg is nil
func NewObserverHandler(handler slog.Handler, observer RecordCollector) (ObserverHandler, error) {
	var err error
	if handler == nil {
		err = errors.New("handler passed to ObserverHandler constructor can not be nil")
	} else if observer == nil {
		err = errors.New("observer passed to ObserverHandler constructor can not be nil")
	}
	return ObserverHandler{handler, observer}, err
}

// ObserverHandler constructor with default handler
func NewDefaultObserverHandler(observer RecordCollector) (ObserverHandler, error) {
	h := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})
	return NewObserverHandler(h, observer)
}

func (h ObserverHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h ObserverHandler) Handle(ctx context.Context, record slog.Record) error {
	if err := h.Handler.Handle(ctx, record); err != nil {
		return err
	}
	h.RecordCollector.addRecord(record)
	return nil
}

func (h ObserverHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.Handler = h.Handler.WithAttrs(attrs)
	return h
}

func (h ObserverHandler) WithGroup(name string) slog.Handler {
	h.Handler = h.Handler.WithGroup(name)
	return h
}
