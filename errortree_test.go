package errortree

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFlatten(t *testing.T) {
	// A non-Tree error should return nil
	require.Nil(t, Flatten(errors.New("test")))

	// Simple, single level tree
	tree := &Tree{
		Delimiter: ".",
		Errors: map[string]error{
			"a": errors.New("test0"),
			"b": errors.New("test1"),
			"c": errors.New("test2"),
		},
	}

	flattened := Flatten(tree)
	require.NotNil(t, flattened)
	require.Len(t, flattened, 3)
	require.EqualValues(t, tree.Errors, flattened)

	// Multi-level tree
	tree = &Tree{
		Delimiter: ".",
		Errors: map[string]error{
			"a": errors.New("test0"),
			"b": &Tree{
				Delimiter: ":",
				Errors: map[string]error{
					"ba": errors.New("test1"),
					"bb": errors.New("test2"),
				},
			},
			"c": &Tree{
				Delimiter: "!",
				Errors: map[string]error{
					"ca": &Tree{
						Delimiter: "-",
						Errors: map[string]error{
							"caa": errors.New("test3"),
							"cab": errors.New("test4"),
						},
					},
					"cb": errors.New("test5"),
				},
			},
		},
	}

	expected := map[string]error{
		"a":        errors.New("test0"),
		"b.ba":     errors.New("test1"),
		"b.bb":     errors.New("test2"),
		"c.ca.caa": errors.New("test3"),
		"c.ca.cab": errors.New("test4"),
		"c.cb":     errors.New("test5"),
	}

	flattened = Flatten(tree)
	require.NotNil(t, flattened)
	require.Len(t, flattened, 6)
	require.EqualValues(t, expected, flattened)

	// Simple recursion (single level)
	tree = &Tree{
		Delimiter: ".",
		Errors: map[string]error{
			"a": errors.New("test0"),
		},
	}
	tree.Errors["b"] = tree
	expected = map[string]error{
		"a": errors.New("test0"),
	}

	flattened = Flatten(tree)
	require.NotNil(t, flattened)
	require.Len(t, flattened, 1)
	require.EqualValues(t, expected, flattened)

	// Multi-level recursion
	tree = &Tree{
		Delimiter: ".",
		Errors: map[string]error{
			"a": errors.New("test0"),
		},
	}
	childTree := &Tree{
		Delimiter: ".",
		Errors:    make(map[string]error, 1),
	}
	childTree.Errors["c"] = tree
	tree.Errors["b"] = childTree
	expected = map[string]error{
		"a": errors.New("test0"),
	}

	flattened = Flatten(tree)
	require.NotNil(t, flattened)
	require.Len(t, flattened, 1)
	require.EqualValues(t, expected, flattened)
}

func TestGet(t *testing.T) {
	// Non-tree get should return nil
	require.Nil(t, Get(errors.New("test"), "test"))

	tree := &Tree{
		Delimiter: ".",
		Errors: map[string]error{
			"a": errors.New("test0"),
			"c": &Tree{
				Errors: map[string]error{
					"a": errors.New("test1"),
				},
			},
		},
	}

	// Top-level: non-existing key
	require.Nil(t, Get(tree, "b"))

	// Top-level: existing key
	require.EqualError(t, Get(tree, "a"), "test0")

	// Nested: non-existing key
	require.Nil(t, Get(tree, "c", "b"))

	// Nested: existing key
	require.EqualError(t, Get(tree, "c", "a"), "test1")

	// Nested: non-tree key
	require.NoError(t, Get(tree, "a", "b"))
}

func TestGetAny(t *testing.T) {
	// Non-tree get should return the error
	err := errors.New("test")
	require.EqualError(t, GetAny(err, "test"), err.Error())

	tree := &Tree{
		Delimiter: ".",
		Errors: map[string]error{
			"a": errors.New("test0"),
			"c": &Tree{
				Errors: map[string]error{
					"a": errors.New("test1"),
				},
			},
		},
	}

	// Top-level: non-existing key
	require.Nil(t, GetAny(tree, "b"))

	// Top-level: existing key
	require.EqualError(t, GetAny(tree, "a"), "test0")

	// Nested: non-existing key
	require.EqualError(t, GetAny(tree, "c", "b"), tree.Errors["c"].Error())

	// Nested: existing key
	require.EqualError(t, GetAny(tree, "c", "a"), "test1")

	// Nested: non-tree key
	require.EqualError(t, GetAny(tree, "a", "b"), "test0")
}

func TestKeys(t *testing.T) {
	// Non-tree should return nil
	require.Nil(t, Keys(errors.New("test")))

	// Simple tree
	tree := &Tree{
		Delimiter: ".",
		Errors: map[string]error{
			"a": errors.New("test0"),
			"b": errors.New("test1"),
		},
	}
	require.EqualValues(t, []string{"a", "b"}, Keys(tree))

	// Nested trees
	tree = &Tree{
		Delimiter: ".",
		Errors: map[string]error{
			"a": errors.New("test0"),
			"b": errors.New("test1"),
			"c": &Tree{
				Delimiter: "/",
				Errors: map[string]error{
					"a": errors.New("test2"),
					"b": errors.New("test3"),
				},
			},
		},
	}
	require.EqualValues(t, []string{"a", "b", "c.a", "c.b"}, Keys(tree))
}

func TestSet(t *testing.T) {
	// Set from nil
	tree := Set(nil, "a", errors.New("test")).(*Tree)
	require.NotNil(t, tree)
	require.EqualValues(t, tree.Errors, map[string]error{
		"a": errors.New("test"),
	})
	require.NotNil(t, tree.Formatter)
	require.EqualValues(t, DefaultDelimiter, tree.Delimiter)

	// Set with nil error: should be a no-op
	tree = Set(tree, "b", nil).(*Tree)
	require.NotNil(t, tree)
	keyValue, keyExists := tree.Errors["b"]
	require.False(t, keyExists, "key expected to not exists, but key exists with value %+v", keyValue)

	// Set from existing tree
	tree2 := Set(tree, "b", errors.New("test2")).(*Tree)
	require.NotNil(t, tree2)
	require.EqualValues(t, tree, tree2)
	require.EqualValues(t, tree2.Errors, map[string]error{
		"a": errors.New("test"),
		"b": errors.New("test2"),
	})

	// Set nested
	tree3 := Set(tree, "c", Set(nil, "a", errors.New("test3"))).(*Tree)
	require.NotNil(t, tree3)
	require.Len(t, tree3.Errors, 3)
	nested := tree3.Errors["c"]
	require.NotNil(t, nested)
	require.IsType(t, &Tree{}, nested)
	require.EqualError(t, nested.(*Tree).Errors["a"], "test3")

	// Set on non-tree: should panic
	func() {
		defer func() {
			r := recover()
			require.NotNil(t, r)
			require.EqualValues(t, r, "Cannot set error: not an *errortree.Tree.")
		}()

		Set(errors.New("test"), "a", errors.New("test0"))
	}()
}

func TestAdd(t *testing.T) {
	// Add from nil
	tree := Add(nil, "a", errors.New("test")).(*Tree)
	require.NotNil(t, tree)
	require.EqualValues(t, tree.Errors, map[string]error{
		"a": errors.New("test"),
	})
	require.NotNil(t, tree.Formatter)
	require.EqualValues(t, DefaultDelimiter, tree.Delimiter)

	// Add from existing tree
	tree2 := Add(tree, "b", errors.New("test2")).(*Tree)
	require.NotNil(t, tree2)
	require.EqualValues(t, tree, tree2)
	require.EqualValues(t, tree2.Errors, map[string]error{
		"a": errors.New("test"),
		"b": errors.New("test2"),
	})

	// Add nested
	tree3 := Add(tree, "c", Set(nil, "a", errors.New("test3"))).(*Tree)
	require.NotNil(t, tree3)
	require.Len(t, tree3.Errors, 3)
	nested := tree3.Errors["c"]
	require.NotNil(t, nested)
	require.IsType(t, &Tree{}, nested)
	require.EqualError(t, nested.(*Tree).Errors["a"], "test3")

	// Add on non-tree: should panic
	func() {
		defer func() {
			r := recover()
			require.NotNil(t, r)
			require.EqualValues(t, r, "Cannot set error: not an *errortree.Tree.")
		}()

		Set(errors.New("test"), "a", errors.New("test0"))
	}()

	// Add on tree with existing key: should panic
	func() {
		defer func() {
			r := recover()
			require.NotNil(t, r)
			require.EqualValues(t, r, "Cannot add error: key a exists.")
		}()

		Add(tree3, "a", errors.New("test0"))
	}()
}
