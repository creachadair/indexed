// Program mktype generates filter functions for specified types.  It is
// intended for use with "go generate".
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
)

var (
	typeName = flag.String("type", "", "Generated type name")
	baseType = flag.String("base", "", "Base type")
	funcName = flag.String("func", "", "Filter function type name")
	pkgName  = flag.String("pkg", "", "Package name (optional)")
	outputTo = flag.String("out", "", "Output file path (optional)")
	doAppend = flag.Bool("append", false, "Append output rather than overwriting")
)

func main() {
	flag.Parse()
	switch {
	case *typeName == "":
		log.Fatal("You must provide a non-empty --type name")
	case *baseType == "":
		log.Fatal("You must provide a non-empty --base type")
	case *funcName == "":
		log.Fatal("You must provide a non-empty --func type name")
	case *pkgName == "" && !*doAppend:
		log.Fatal("You must provide a --pkg name when --append=false")
	}
	out := os.Stdout
	if *outputTo != "" {
		flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
		if *doAppend {
			flags &^= os.O_TRUNC
			flags |= os.O_APPEND
		}

		f, err := os.OpenFile(*outputTo, flags, 0644)
		if err != nil {
			log.Fatalf("Unable to append to %q: %v", *outputTo, err)
		}
		out = f
	}

	buf := bytes.NewBuffer(nil)

	// If specified, emit a package name.
	if *pkgName != "" {
		fmt.Fprintf(buf, "package %s\n", *pkgName)
		fmt.Fprint(buf, "// Generated code, do not edit (see gentypes.go).\n")
	}

	// Generate the base type definition.
	fmt.Fprintf(buf, "\ntype %[1]s []%[2]s\n", *typeName, *baseType)

	// Generate the required methods for the interface.
	fmt.Fprintf(buf, "func (t %[1]s) Len() int { return len(t) }\n", *typeName)
	fmt.Fprintf(buf, "func (t %[1]s) Swap(i, j int) { t[i], t[j] = t[j], t[i] }\n", *typeName)

	// Generate a constructor function.
	fmt.Fprintf(buf, filterFunc, *funcName, *baseType, *typeName)

	code, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalf("Error in generated code: %v\n%s", err, buf.String())
	}
	fmt.Fprintln(out, string(code))
	if err := out.Close(); err != nil {
		log.Printf("Warning: error closing output: %v", err)
	}
}

const filterFunc = `
// Filter%[1]s modifies *ss in-place to remove any elements for which keep returns
// false. Relative input order is preserved. If ss == nil, this function panics.
func Filter%[1]s(ss *[]%[2]s, keep func(%[2]s) bool) {
   *ss = (*ss)[:Partition(%[3]s(*ss), func(i int) bool {
      return keep((*ss)[i])
   })]
}
`
