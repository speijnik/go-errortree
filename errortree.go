// Package errortree provides primitives for working with errors in tree structure
//
// errortree is intended to be used in places where errors are generated
// from an arbitrary tree structure, like the validation of a configuration
// file.
// This allows adding additional context as to why an error has happened
// in a clean and structured way.
//
// errortree fully supports nesting of multiple trees, including simplified
// retrieval of errors which, among other things, should help remove repeated boilerplate code
// from unit tests.
package errortree

import (
	"sort"
)

// DefaultDelimiter defines the delimiter that is by default used
// for representing paths to nested Errors
const DefaultDelimiter = ":"

// Keys returns all error keys present in a given tree.
//
// The value returned by this function is an alphabetically sorted, flattened
// list of all keys in a tree of Tree structs.
//
// The delimiter configured for the top-level tree is guaranteed to be used
// throughout the complete tree.
func Keys(err error) []string {
	tree, isTree := GetTree(err)

	// If the supplied error is not a *Tree we return nil
	if !isTree {
		return nil
	}

	flattened := flatten(tree, tree.getDelimiter(), nil)
	keys := make([]string, 0, len(flattened))
	for key := range flattened {
		keys = append(keys, key)
	}

	// Sort the keys before returning them
	sort.Strings(keys)

	return keys
}

func set(tree *Tree, key string, err error) *Tree {
	if err == nil {
		return tree
	}

	if tree == nil {
		tree = New()
	}

	errors := tree.getErrors()
	errors[key] = err

	return tree
}

// Set creates or replaces an error under a given key in a tree.
//
// The parent value may be nil, in which case a new *Tree is created, to which the
// key is added and the new *Tree is returned.
// Otherwise the *Tree to which the key was added is returned.
func Set(parent error, key string, err error) error {
	tree, isTree := GetTree(parent)

	// Add only works on an *errortree.Tree, panic if we received another error
	if parent != nil && !isTree {
		panic("Cannot set error: not an *errortree.Tree.")
	}

	if tree = set(tree, key, err); tree != nil {
		return tree
	}
	return nil
}

// Add adds an error under a given key to the provided tree.
//
// This function panics if the key is already present in the tree.
// Otherwise it behaves like Set.
func Add(parent error, key string, err error) error {
	tree, isTree := GetTree(parent)
	if tree != nil && isTree {
		if _, keyExists := tree.getErrors()[key]; keyExists {
			panic("Cannot add error: key " + key + " exists.")
		}
	}

	if tree = set(tree, key, err); tree != nil {
		return tree
	}
	return nil
}

// Get retrieves the error for the given key from the provided error.
// The path parameter may be used for specifying a nested error's key.
//
// If the error is not an errortree.Tree or the child cannot be found on the exact
// path this function returns nil.
func Get(err error, key string, path ...string) error {
	tree, isTree := GetTree(err)
	if !isTree {
		return nil
	}

	return get(tree, false, key, path...)
}

// GetAny retrieves the error for a given key from the tree.
// The path parameter may be used for specifying a nested error's key.
//
// This function returns the most-specific match:
//
// If the provided error is not an errortree.Tree, the provided error is returned.
// If at any step the path cannot be fully followed, the previous error on the path will be returned.
func GetAny(err error, key string, path ...string) error {
	tree, isTree := GetTree(err)
	if !isTree {
		return err
	}

	return get(tree, true, key, path...)
}

func get(tree *Tree, returnAnyChild bool, key string, path ...string) error {
	child, keyExists := tree.getErrors()[key]
	if !keyExists {
		return nil
	} else if len(path) == 0 {
		return child
	}

	childTree, isTree := GetTree(child)
	if !isTree && returnAnyChild {
		return child
	} else if !isTree {
		return nil
	}

	childErr := get(childTree, returnAnyChild, path[0], path[1:]...)

	if childErr == nil && returnAnyChild {
		return child
	}

	return childErr
}

// Flatten returns the error tree in flattened form.
//
// Each error inside the complete tree is stored under its full key.
// The full key is constructed from the each error's path inside the tree
// and joined together with the tree's delimiter.
func Flatten(err error) map[string]error {
	tree, isTree := GetTree(err)
	if !isTree {
		return nil
	}

	return flatten(tree, tree.getDelimiter(), nil)
}

func flatten(tree *Tree, delimiter string, visited []*Tree, keyPrefix ...string) map[string]error {
	for _, visitedTree := range visited {
		if tree == visitedTree {
			return map[string]error{}
		}
	}
	visited = append(visited, tree)

	errors := tree.getErrors()
	errorMap := make(map[string]error, len(errors))

	for key, err := range errors {
		if childTree, isTree := GetTree(err); isTree {
			childPrefix := append(keyPrefix, key)
			for childKey, childErr := range flatten(childTree, delimiter, visited, childPrefix...) {
				errorMap[key+delimiter+childKey] = childErr
			}
		} else {
			errorMap[key] = err
		}
	}

	return errorMap
}
