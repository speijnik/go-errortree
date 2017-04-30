package errortree

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimpleFormatter(t *testing.T) {
	// Test single error
	require.EqualValues(t, "1 error occurred:\n\n* key: value", SimpleFormatter(
		map[string]error{
			"key": errors.New("value"),
		},
	))

	// Test multiple Errors, including alphabetical order
	// Use error values that would break the correct order to ensure we are ordering by key and not by value
	require.EqualValues(t, "3 errors occurred:\n\n* a: c\n* b: a\n* c: b", SimpleFormatter(
		map[string]error{
			"a": errors.New("c"),
			"b": errors.New("a"),
			"c": errors.New("b"),
		},
	))
}
