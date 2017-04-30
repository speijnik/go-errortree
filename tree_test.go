package errortree

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	// Validate default options
	tree := New()
	require.NotNil(t, tree)
	require.EqualValues(t, DefaultDelimiter, tree.Delimiter)
	require.NotNil(t, tree.Formatter)
	require.NotNil(t, tree.Errors)

	// Misc: validate that getErrors (internal function) initializes the Errors map
	tree = &Tree{}
	require.NotNil(t, tree.getErrors())
	require.NotNil(t, tree.Errors)
}

func TestTree_ErrorOrNil(t *testing.T) {
	// Create new error
	tree := &Tree{}
	require.NotNil(t, tree)
	require.Nil(t, tree.ErrorOrNil())

	// Set t to nil, this should not panic but rather return nil
	tree = nil
	require.Nil(t, tree.ErrorOrNil())

	tree = &Tree{
		Errors: map[string]error{
			"test": errors.New("test"),
		},
	}
	require.NotNil(t, tree.ErrorOrNil())
}

func TestTree_WrappedErrors(t *testing.T) {
	tree := &Tree{
		// Use keys and error messages that would be sorted differently.
		// This is required to ensure that sorting is done using the keys.
		Errors: map[string]error{
			"a": errors.New("c"),
			"b": errors.New("a"),
			"c": errors.New("b"),
		},
	}

	wrappedErrors := tree.WrappedErrors()
	require.Len(t, wrappedErrors, 3)
	require.EqualValues(t, wrappedErrors, []error{
		errors.New("c"),
		errors.New("a"),
		errors.New("b"),
	})
}

func TestTree_Error(t *testing.T) {
	treeErrors := map[string]error{
		"a": errors.New("c"),
		"b": &Tree{
			Errors: map[string]error{
				"0": errors.New("test0"),
				"1": errors.New("test1"),
			},
		},
		"c": errors.New("b"),
	}

	flattenedErrors := map[string]error{
		"a":   errors.New("c"),
		"b.0": errors.New("test0"),
		"b.1": errors.New("test1"),
		"c":   errors.New("b"),
	}

	formatter := func(errorMap map[string]error) string {
		// Validate that we receive the expected flattened error map into the formatter
		require.EqualValues(t, flattenedErrors, errorMap)
		return "formatter_called"
	}

	tree := &Tree{
		Formatter: formatter,
		Errors:    treeErrors,
		Delimiter: ".",
	}

	// Call the formatter, but only require that formatter_called is returned
	require.EqualValues(t, "formatter_called", tree.Error())

}

func TestGetTree(t *testing.T) {
	err := errors.New("Test")
	tree, isTree := GetTree(err)
	require.Nil(t, tree)
	require.False(t, isTree)

	err = &Tree{}
	tree, isTree = GetTree(err)
	require.EqualValues(t, err, tree)
	require.True(t, isTree)
}
