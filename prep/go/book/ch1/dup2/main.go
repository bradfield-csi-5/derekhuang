// Dup2 prints the count and text of lines that appear more than once
// in the input.  It reads from stdin or from a list of named files.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	counts := make(map[string]int)
	origins := make(map[string][]string)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts, origins)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts, origins)
			f.Close()
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\t%s\n", n, line, strings.Join(origins[line], ", "))
		}
	}
}

func countLines(f *os.File, counts map[string]int, origins map[string][]string) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		line := input.Text()
		if strings.TrimSpace(line) != "" {
			counts[line]++
			found := false
			for i := 0; i < len(origins[line]); i++ {
				if origins[line][i] == f.Name() {
					found = true
					break
				}
			}
			if !found {
				origins[line] = append(origins[line], f.Name())
			}
		}
	}
	// NOTE: ignoring potential errors from input.Err()
}
