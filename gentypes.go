package indexed

// Instructions to "go generate" to produce wrapper types that implement the
// Filterable interface. To add additional wrappers, append them to this list
// and re-run "go generate".

//go:generate go run mktype/mktype.go -type stringSwapper -base string -func Strings -out types.go -pkg indexed
//go:generate go run mktype/mktype.go -type intSwapper -base int -func Ints -out types.go -append
