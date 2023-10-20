package main

import (
	"bytes"
	"log"
	"os"
	"sort"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

const src string = `package foo

import (
	"fmt"
	"time"
)

func baz() {
	fmt.Println("Hello, world!")
}

type A int

const b = "testing"

func bar() {
	fmt.Println(time.Now())
}`

// Moves all top-level functions to the end, sorted in alphabetical order.
// The "source file" is given as a string (rather than e.g. a filename).
func SortFunctions(src string) (string, error) {
	f, err := decorator.Parse(src)
	if err != nil {
		return "", err
	}
	sort.Slice(f.Decls, func(i, j int) bool {
		x, xok := f.Decls[i].(*dst.FuncDecl)
		y, yok := f.Decls[j].(*dst.FuncDecl)
		if xok && !yok {
			return false
		} else if !xok && yok {
			return true
		} else if !xok && !yok {
			return false
		}
		return x.Name.Name < y.Name.Name
	})
	var buf bytes.Buffer
	err = decorator.Fprint(&buf, f)
	return buf.String(), err
}

func main() {
	f, err := decorator.Parse(src)
	if err != nil {
		log.Fatal(err)
	}

	// Print AST
	err = dst.Fprint(os.Stdout, f, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Convert AST back to source
	err = decorator.Print(f)
	if err != nil {
		log.Fatal(err)
	}
}
