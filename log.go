package slogt

import (
	"strings"
	"time"

	"golang.org/x/exp/slog"
)

// Wrap a slog.Record & provides getters on it
type Log struct {
	// embedding would expose record fields
	// and we don't want to allow modification
	record slog.Record
	// store record attributes for reuse in FindAttribute() method
	_attributes []slog.Attr
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

// find an attribute by its key (slog.Attr.Key)
// for nested attributes (slog.Group) you may use dot notation e.g. user.profile.age
// when not found, a zero-ed attribute is returned
func (l Log) FindAttribute(key string) (attribute slog.Attr, found bool) {

	attributes := l.getAttributes()
	subkeys := strings.Split(key, ".")

	// each loop checks attributes for a match with the current subkey
	// if an attribute is found and is a group, its value will replace attributes
	// and the current subkey is remove from subkeys
ext:
	for {

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

// extract attributes from slog.Record once
func (l Log) getAttributes() []slog.Attr {

	if len(l._attributes) == 0 {

		l.record.Attrs(func(attr slog.Attr) bool {
			l._attributes = append(l._attributes, attr)
			return true
		})
	}

	return l._attributes
}
