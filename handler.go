package slogt

import (
	"context"
	"errors"
	"io"

	"golang.org/x/exp/slog"
)

/*
This is a custom slog.Handler.

It is composed of a common Handler and an Observer.
The common Handler helps implementing the slog.Handler interface without entirely rewriting its methods.
The Observer will store each slog.Record passed by the Logger (see Handle()).

The handler groups and attributes are not exposed by the common Handler.
We need to capture and store them (see WithGroup() and WithAttrs()).
*/

// An observer must implement this interface to be accepted by the ObserverHandler.
// see ObserverHandler.Handle()
type HandlerObserver interface {
	addLog(log Log)
}

type ObserverHandler struct {
	handler  slog.Handler
	observer HandlerObserver
	// groups and shared attributes are not exposed by the common Handler
	// we have to capture and store them here
	groups []string
	attrs  []slog.Attr
}

// Constructor
// An error is returned if the handler or observer arg is nil.
func NewObserverHandler(observer HandlerObserver, handler slog.Handler) (ObserverHandler, error) {
	var err error
	if handler == nil {
		err = errors.New("handler passed to ObserverHandler constructor can not be nil")
	} else if observer == nil {
		err = errors.New("observer passed to ObserverHandler constructor can not be nil")
	}
	return ObserverHandler{
		handler:  handler,
		observer: observer,
		groups:   []string{},
		attrs:    []slog.Attr{},
	}, err
}

// Constructor without handler argument.
// An error is returned if the observer arg is nil.
// Don't expect log output, the handler uses a io.Discard Writer.
func NewDefaultObserverHandler(observer HandlerObserver) (ObserverHandler, error) {
	handler := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})
	return NewObserverHandler(observer, handler)
}

func (oh ObserverHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return oh.handler.Enabled(ctx, level)
}

func (oh ObserverHandler) Handle(ctx context.Context, record slog.Record) error {
	if err := oh.handler.Handle(ctx, record); err != nil {
		return err
	}
	log := NewLog(record, oh)
	oh.observer.addLog(log)
	return nil
}

// Get a new ObserverHandler (type assertion) with additional shared attributes.
func (oh ObserverHandler) WithAttrs(attrs []slog.Attr) slog.Handler {

	newHandler := oh.handler.WithAttrs(attrs)

	return ObserverHandler{
		handler:  newHandler,
		observer: oh.observer, // keep same observer with stored logs
		groups:   oh.groups,
		attrs:    append(oh.attrs, attrs...),
	}
}

// Get a new ObserverHandler (type assertion) with an additional group.
func (oh ObserverHandler) WithGroup(name string) slog.Handler {

	// see slog.handler.go: func (h *commonHandler) withGroup(name string)
	if name == "" {
		return oh
	}

	newHandler := oh.handler.WithGroup(name)

	return ObserverHandler{
		handler:  newHandler,
		observer: oh.observer, // keep same observer with stored logs
		groups:   append(oh.groups, name),
		attrs:    oh.attrs,
	}
}

// Get this handlers's groups.
func (oh ObserverHandler) Groups() []string {
	return oh.groups
}

// Get this handlers's attributes.
func (oh ObserverHandler) Atttributes() []slog.Attr {
	return oh.attrs
}

// Search fr an attribute by its Key in this handler's attributes.
// If not found, the returned attribute has a nil Value.
func (oh ObserverHandler) FindAttribute(key string) (attribute slog.Attr, found bool) {
	return findAttribute(key, oh.attrs)
}
