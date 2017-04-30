package errortree_test

import (
	"errors"
	"fmt"
	"strings"

	"github.com/speijnik/go-errortree"
)

func ExampleAdd() {
	var err error

	// Using Add on a nil-error automatically creates an error tree and adds the desired error
	err = errortree.Add(err, "test", errors.New("test error"))
	fmt.Println(err.Error())
	// Output: 1 error occurred:
	//
	// * test: test error

}

func ExampleAdd_nested() {
	// Create an error which will acts as the child for our top-level error
	childError := errortree.Add(nil, "test0", errors.New("child error"))
	// Add another error to our child
	childError = errortree.Add(childError, "test1", errors.New("another child error"))

	// Create the top-level error, adding the child in the process
	err := errortree.Add(nil, "child", childError)
	// Add another top-level error
	err = errortree.Add(err, "second", errors.New("top-level error"))

	fmt.Println(err.Error())
	// Output: 3 errors occurred:
	//
	// * child:test0: child error
	// * child:test1: another child error
	// * second: top-level error
}

func ExampleAdd_duplicate() {
	// Add an error with the key "test"
	err := errortree.Add(nil, "test", errors.New("test error"))

	// Recover from the panic, this will the output we expect below
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	err = errortree.Add(err, "test", errors.New("key re-used"))
	// Output: Cannot add error: key test exists.
}

func ExampleSet() {
	var err error

	// Using Set on a nil-error automatically creates an error tree and adds the desired error
	err = errortree.Set(err, "test", errors.New("test error"))
	fmt.Println(err.Error())
	// Output: 1 error occurred:
	//
	// * test: test error
}

func ExampleSet_duplicate() {
	// Add an error with the key "test"
	err := errortree.Add(nil, "test", errors.New("test error"))

	// Call Set on the key, which will override it
	err = errortree.Set(err, "test", errors.New("key re-used"))
	fmt.Println(err.Error())
	// Output: 1 error occurred:
	//
	// * test: key re-used
}

func ExampleFlatten() {
	tree := &errortree.Tree{
		Errors: map[string]error{
			"a": errors.New("top-level"),
			"b": &errortree.Tree{
				Errors: map[string]error{
					"c": errors.New("nested"),
				},
			},
		},
	}

	flattened := errortree.Flatten(tree)
	// Sort keys alphabetically so we get reproducible output
	keys := errortree.Keys(tree)
	for _, key := range keys {
		fmt.Println("key: " + key + ", value: " + flattened[key].Error())
	}

	// Output: key: a, value: top-level
	// key: b:c, value: nested
}

func ExampleKeys() {
	tree := &errortree.Tree{
		Errors: map[string]error{
			"a":    errors.New("top-level"),
			"test": errors.New("test"),
			"b": &errortree.Tree{
				Errors: map[string]error{
					"c": errors.New("nested"),
				},
			},
		},
	}

	fmt.Println(strings.Join(errortree.Keys(tree), ", "))
	// Output: a, b:c, test
}

func ExampleGet() {
	tree := &errortree.Tree{
		Errors: map[string]error{
			"a":    errors.New("top-level"),
			"test": errors.New("test"),
			"b": &errortree.Tree{
				Errors: map[string]error{
					"c": errors.New("nested"),
				},
			},
		},
	}

	// Get can be used to retrieve an error by its key
	fmt.Println(errortree.Get(tree, "a"))

	// Nested retrieval is supported as well
	fmt.Println(errortree.Get(tree, "b", "c"))
	// Output: top-level
	// nested
}

func ExampleGet_nested() {
	tree := &errortree.Tree{
		Errors: map[string]error{
			"a":    errors.New("top-level"),
			"test": errors.New("test"),
			"b": &errortree.Tree{
				Errors: map[string]error{
					"c": errors.New("nested"),
				},
			},
		},
	}

	// Get tries to resolve the path exactly and returns nil if
	// the path does not exist
	fmt.Println(errortree.Get(tree, "b", "non-existent"))
	// Output: <nil>
}

func ExampleGet_non_tree() {
	// Get returns nil if the passed error is not a tree

	fmt.Println(errortree.Get(errors.New("test"), "key"))
	// Output: <nil>
}

func ExampleGetAny_non_tree() {
	// GetAny always returns the error it got passed even if it is not a tree

	fmt.Println(errortree.GetAny(errors.New("test"), "key"))
	// Output: test
}

func ExampleGetAny_nested() {
	// When GetAny does not encounter an exact match in the tree, it returns the most-specific match

	tree := &errortree.Tree{
		Errors: map[string]error{
			"a":    errors.New("top-level"),
			"test": errors.New("test"),
			"b": &errortree.Tree{
				Errors: map[string]error{
					"c": errors.New("nested"),
				},
			},
		},
	}

	// Get tries to resolve the path exactly and returns nil if
	// the path does not exist
	fmt.Println(errortree.GetAny(tree, "b", "non-existent"))
	// Output: 1 error occurred:
	//
	// * c: nested
}
