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

// search for a attribute shared by all logs (for a given logger and handler)
// shared attributes are defined by a slog.Logger or a slog.Handler
func (l Log) FindSharedAttribute(key string) (attribute slog.Attr, found bool) {
	return findAttribute(key, l.sharedAttributes)
}

func (l Log) FindBuiltInAttribute(key string) (attribute slog.Attr, found bool) {
	return findAttribute(l.stripGroupNames(key), l.GetBuiltInAttributes())
}

// return built-in attributes
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

// return group names defined by a logger and/or a handler
// names are separated by a dot e.g. "auth.admin"
func (l Log) GroupNames() string {
	return strings.Join(l.groups, ".")
}

// find an attribute by its Key in attributes
func findAttribute(key string, attributes []slog.Attr) (attribute slog.Attr, found bool) {

	subkeys := splitKey(key)

ext:
	for {
		// search in attributes for a match with first item in subkeys
		// if found, subkeys is shifted to the left and the attribute Value - if a group - becomes attributes

		subkey := subkeys[0]
		subkeyAttrFound := false
		isLastSubkey := len(subkeys) == 1

		for _, attr := range attributes {

			if attr.Key == subkey {

				subkeyAttrFound = true

				if isLastSubkey {
					found = true
					attribute = attr
					break ext
				}

				// make sure found attribute is a group to continue
				if attr.Value.Kind() == slog.KindGroup {
					// update attributes and subkeys for next loop
					attributes = attr.Value.Group()
					subkeys = subkeys[1:]
				} else {
					break ext
				}
			}
		}

		if !subkeyAttrFound {
			break
		}

	}

	return attribute, found
}

// remove group names (joined by a dot) at the start of the given key
// e.g. key = "app1.user.id" ang groups = ["app1"] will return "user.id"
func (l Log) stripGroupNames(key string) string {

	groups := joinKeys(l.groups)
	index := strings.Index(key, groups)
	if index == 0 {
		// add one to remove the dot between app1 and user
		key = key[len(groups)+1:]
	}

	return key
}

// split an attribute key with "."
func splitKey(key string) []string {
	return strings.Split(key, ".")
}

// join attribute subkeys with "."
func joinKeys(keys []string) string {
	return strings.Join(keys, ".")
}
