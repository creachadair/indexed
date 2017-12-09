# filter

This repository provides a package that can partition the contents of an
arbitrary indexed collection (typically a slice), in-place according to a
selection rule.  After partitioning, all the selected elements are at low-order
indices of the collection, in their original relative order; the unselected
elements follow in arbitrary order. This operation takes linear time in the
size of the collection, and constant space overhead for bookkeeping.

This operation permits a collection to be filtered efficiently without
redundant copying or movement of data. Some convenience functions are provided
for applying this to common built-in data types.

View documentation on [GoDoc](http://godoc.org/bitbucket.org/creachadair/filter).
