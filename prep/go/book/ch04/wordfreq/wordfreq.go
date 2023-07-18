/*
Exercise 4.9: Write a program wordfreq to report the frequency of each word in an input text file.
Call input.Split(bufio.ScanWords) before the first call to Scan to break the input into words instead of lines.
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

func main() {
	var words []string
	files := os.Args[1:]

	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "wordfreq: a file name is required")
		return
	}

	counts := make(map[string]int)
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "wordfreq: error opening file "+file)
		}
		input := bufio.NewScanner(f)
		input.Split(bufio.ScanWords)
		for input.Scan() {
			counts[format(input.Text())]++
		}
	}

	// sort to print in alphabetical order
	for word := range counts {
		words = append(words, word)
	}
	sort.Strings(words)
	fmt.Printf("Word\tcount\n")
	for _, w := range words {
		p := "%q\t%d\n"

		// line up counts with two tabs if the word is short
		if len(w) < 6 {
			p = "%q\t\t%d\n"
		}
		fmt.Printf(p, w, counts[w])
	}
}

func format(s string) string {
	// remove punctuation and make lowercase
	return regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(strings.ToLower(s), "")
}
