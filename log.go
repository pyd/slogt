package slogt

import (
	"strings"
	"time"

	"golang.org/x/exp/slog"
)

/*
Log contains data related to a single log:
  - time
  - level
  - message
  - built-in attributes i.e. defined in a slog.Record
    such attributes are passed as args of a Logger method like Info(), Error(), Log(), LogAttrs()...
    or added to an existing record via Record.Add() or Record.AddAttrs()
  - handler also provides access to its groups and attributes
*/
type Log struct {
	// embedding would expose record fields
	// and we don't want to allow modifications
	record slog.Record
	// provides access to handler groups and attributes
	handler ObserverHandler
	// built-in attributes, stored once to be reused by FindAttribute()
	_builtInAttributes []slog.Attr
}

func NewLog(record slog.Record, handler ObserverHandler) Log {
	return Log{
		record:  record,
		handler: handler,
		// // groups:           groups,
		// sharedAttributes: attrs,
	}
}

func (l Log) Message() string {
	return l.record.Message
}

func (l Log) Time() time.Time {
	return l.record.Time
}

func (l Log) Level() slog.Level {
	return l.record.Level
}

// Search for a shared attribute (defined at logger or handler level) by its key.
// Note: shared attributes are not prefixed with handler groups.
func (l Log) FindSharedAttribute(key string) (attribute slog.Attr, found bool) {
	return l.handler.FindAttribute(key)
}

// Search for a built-in attribute (defined in Logger methods e.g. Warn()) by its key.
// The key can be prefixed with handler groups but it is not a requirement.
func (l Log) FindBuiltInAttribute(key string) (attribute slog.Attr, found bool) {
	return findAttribute(l.stripGroupNames(key), l.GetBuiltInAttributes())
}

// Get built-in attributes.
// Note: attributes keys are not prefixed with handler groups.
func (l Log) GetBuiltInAttributes() []slog.Attr {

	// extract built-in attributes from record once
	// and store it for reuse
	if len(l._builtInAttributes) == 0 {

		l.record.Attrs(func(attr slog.Attr) bool {
			l._builtInAttributes = append(l._builtInAttributes, attr)
			return true
		})
	}

	return l._builtInAttributes
}

// Get groups from the handler for this log.
// Names are separated by a dot e.g. "auth.admin".
func (l Log) GroupNames() string {
	return strings.Join(l.handler.Groups(), ".")
}

// remove group names (joined by a dot) at the start of the given key.
// e.g. key = "app1.user.id" ang groups = ["app1"] will return "user.id".
func (l Log) stripGroupNames(key string) string {

	groups := joinKeys(l.handler.Groups())
	index := strings.Index(key, groups)
	if index == 0 {
		// add one to remove the dot between app1 and user
		key = key[len(groups)+1:]
	}

	return key
}
