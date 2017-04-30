# go-errortree

[![GoDoc](https://godoc.org/github.com/speijnik/go-errortree?status.svg)](https://godoc.org/github.com/speijnik/go-errortree)
[![Build Status](https://travis-ci.org/speijnik/go-errortree.svg)](https://travis-ci.org/speijnik/go-errortree)
[![codecov](https://codecov.io/gh/speijnik/go-errortree/branch/master/graph/badge.svg)](https://codecov.io/gh/speijnik/go-errortree)
[![Go Report Card](https://goreportcard.com/badge/github.com/speijnik/go-errortree)](https://goreportcard.com/report/github.com/speijnik/go-errortree)

`github.com/speijnik/go-errortree` provides functionality for working
with errors structured as a tree.

Structuring errors in such a way may be desired when, for example, validating
structured input such as a configuration file with multiple sections.

A corresponding example can be found in the [example_config_validation_test.go](example_config_validation_test.go) file.

The code is released under the terms of the MIT license.