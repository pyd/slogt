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
  - other attributes ???
  - attributes group names i.e. names defined by a Logger and/or a Handler to qualify (prefix) attributes
*/
type Log struct {
	// embedding would expose record fields
	// and we don't want to allow modifications
	record slog.Record
	// groups qualify/prefix built-in attributes
	// items are joined with a dot
	groups []string
	// attributes set at Logger or Handler level
	// they are added to each log
	sharedAttributes []slog.Attr
	// built-in attributes, stored once to be reused by FindBuiltInAttribute()
	_builtInAttributes []slog.Attr
}

func NewLog(record slog.Record, groups []string, attrs []slog.Attr) Log {
	return Log{
		record:           record,
		groups:           groups,
		sharedAttributes: attrs,
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
	return findAttribute(key, l.sharedAttributes)
}

// Search for a built-in attribute (defined in Logger methods e.g. Warn()) by its key.
// the key can be prefixed with handler groups but it is not a requirement
func (l Log) FindBuiltInAttribute(key string) (attribute slog.Attr, found bool) {
	return findAttribute(l.stripGroupNames(key), l.GetBuiltInAttributes())
}

// Get built-in attributes.
// Note: attributes keys are not prefixed with handler groups
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
	return strings.Join(l.groups, ".")
}
