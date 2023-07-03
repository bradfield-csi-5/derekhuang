/*
Experiment to measure the difference in running time between our potentially inefficient versions and the one that uses strings.Join.

Couldn't figure out why `go test -bench .` didn't work here but it seems like `go test`
didn't like the fact that there were multiple packages in this directory (main, benchmark).

To get it to work, I copied this into a separate directory so it was the only file.
*/
package benchmark

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func BenchmarkFirst(b *testing.B) {
	var args = []string{"foo", "bar", "baz"}
	for i := 0; i < b.N; i++ {
		first(args)
	}
}

func BenchmarkSecond(b *testing.B) {
	var args = []string{"foo", "bar", "baz"}
	for i := 0; i < b.N; i++ {
		second(args)
	}
}

func BenchmarkThird(b *testing.B) {
	var args = []string{"foo", "bar", "baz"}
	for i := 0; i < b.N; i++ {
		third(args)
	}
}

func first(args []string) {
	var s, sep string
	for i := 0; i < len(args); i++ {
		s += sep + args[i]
		sep = " "
	}
	fmt.Println(s)
}

func second(args []string) {
	var s string
	var sep = "\n"
	for i, arg := range args {
		s += strconv.Itoa(i) + " " + arg + sep
	}
	fmt.Println(s)
}

func third(args []string) {
	fmt.Println(strings.Join(args, " "))
}
