/*
Exercise 1.2: Modify the echo program to print the index and value of each of its arguments, one per line.
*/

package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	var s string
	var sep = "\n"
	for i, arg := range os.Args[1:] {
		s += strconv.Itoa(i) + " " + arg + sep
	}
	fmt.Println(s)
}
