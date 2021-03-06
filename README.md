# indexed

[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/creachadair/indexed)
[![Go Report Card](https://goreportcard.com/badge/github.com/creachadair/indexed)](https://goreportcard.com/report/github.com/creachadair/indexed)

The `indexed` package supports sorting and partitioning the contents of an
indexed collection (typically a slice), in-place according to a selection rule.

After partitioning, all the selected elements are at low-order indices of the
collection, in their original relative order; the unselected elements follow in
arbitrary order. This operation takes linear time in the size of the
collection, and constant space overhead for bookkeeping.

This operation permits a collection to be filtered efficiently without
redundant copying or movement of data. Some convenience functions are provided
for applying this to common built-in data types.
