package slogt

import (
	"strings"

	"golang.org/x/exp/slog"
)

// search for an attribute matching the given key
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
