package errortree

import (
	"fmt"
	"sort"
	"strings"
)

// Formatter defines the Formatter type
// This function can expected that the provided map contains a flattened map of all Errors
type Formatter func(map[string]error) string

// SimpleFormatter provides a simple Formatter which returns a message indicating
// how many Errors occurred and details for every error.
// The reported Errors are sorted alphabetically by key.
func SimpleFormatter(errorMap map[string]error) string {
	wrappedErrors := make([]string, len(errorMap))
	pluralSuffix := ""
	if len(errorMap) != 1 {
		pluralSuffix = "s"
	}

	// Sort the keys
	keys := make([]string, 0, len(errorMap))
	for key := range errorMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Construct the individual messages
	for i, key := range keys {
		wrappedErrors[i] = "* " + key + ": " + errorMap[key].Error()
	}

	return fmt.Sprintf("%d error%s occurred:\n\n%s", len(errorMap), pluralSuffix,
		strings.Join(wrappedErrors, "\n"))
}
