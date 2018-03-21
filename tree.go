package errortree

import (
	"sort"
)

var _ error = (*Tree)(nil)

// Tree is an error type which acts as a container for storing
// multiple errors in a tree structure.
type Tree struct {
	// Errors holds the tree's items
	Errors map[string]error
	// Delimiter specifies the tree's delimiter for building nested paths
	Delimiter string
	// Formatter specifies the formatter to use when Error is invoked
	Formatter Formatter
}

func (t *Tree) getErrors() map[string]error {
	if t.Errors == nil {
		t.Errors = make(map[string]error)
	}
	return t.Errors
}

func (t *Tree) getDelimiter() string {
	if t.Delimiter == "" {
		t.Delimiter = DefaultDelimiter
	}

	return t.Delimiter
}

func (t *Tree) getFormatter() Formatter {
	if t.Formatter == nil {
		t.Formatter = SimpleFormatter
	}
	return t.Formatter
}

func (t *Tree) Error() string {
	if t == nil {
		return ""
	}
	formatter := t.getFormatter()

	return formatter(flatten(t, t.getDelimiter(), nil))
}

// ErrorOrNil returns nil if the tree is empty or the tree itself
// otherwise.
func (t *Tree) ErrorOrNil() error {
	if t == nil || len(t.Errors) == 0 {
		return nil
	}
	return t
}

// WrappedErrors returns the errors wrapped by the tree.
//
// The ordering of the returned errors is determined by the alphabetical
// ordering of the corresponding error keys.
func (t *Tree) WrappedErrors() []error {
	errors := t.getErrors()
	wrappedErrors := make([]error, len(errors))
	keys := make([]string, 0, len(errors))
	for key := range errors {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for i, key := range keys {
		wrappedErrors[i] = errors[key]
	}

	return wrappedErrors
}

// New returns a new error tree.
func New() *Tree {
	return &Tree{
		Delimiter: DefaultDelimiter,
		Formatter: SimpleFormatter,
		Errors:    make(map[string]error),
	}
}

// GetTree returns the tree for a given error.
func GetTree(err error) (tree *Tree, isTree bool) {
	tree, isTree = err.(*Tree)
	return
}
