package slogt

import (
	"strings"

	"golang.org/x/exp/slog"
)

// Search for an attribute matching the given key.
// For group attributes, sub keys must be separated by a dot.
func findAttribute(key string, attributes []slog.Attr) (attribute slog.Attr, found bool) {

	subkeys := splitKey(key)

ext:
	for {
		// search attributes for a match with first key in subkeys
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

// Split an attribute key with ".".
func splitKey(key string) []string {
	return strings.Split(key, ".")
}
