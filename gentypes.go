package filter

// Instructions to "go generate" to produce wrapper types that implement the
// Filterable interface. To add additional wrappers, append them to this list
// and re-run "go generate".

//go:generate go run cmd/mktype.go -type stringFilter -base string -cons Strings -out types.go -pkg filter
//go:generate go run cmd/mktype.go -type intFilter -base int -cons Ints -out types.go -append
