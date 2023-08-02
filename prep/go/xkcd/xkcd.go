/*
Exercise 4.12: The popular web comic xkcd has a JSON interface.

For example, a request to https://xkcd.com/571/info.0.json produces a detailed description
of comic 571, one of many favorites. Download each URL (once!) and build an offline index.

Write a tool xkcd that, using this index, prints the URL and transcript of each comic that
matches a search term provided on the command line.
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "[xkcd] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)
var strict = flag.Bool("s", false, "strict search (search terms are ANDed)")

func init() {
	PopulateIndex(logger)
	BuildReverseIndex(logger)
}

func main() {
	flag.Parse()

	var index Index
	b, err := os.ReadFile(IndexFileName)
	if err != nil {
		logger.Fatalln("error reading index:", err)
	}
	err = json.Unmarshal(b, &index)
	if err != nil {
		logger.Fatalln("error unmarshaling index:", err)
	}

	var rindex ReverseIndex
	b, err = os.ReadFile(RevIndexFileName)
	if err != nil {
		logger.Fatalln("error reading reverse index:", err)
	}
	err = json.Unmarshal(b, &rindex)
	if err != nil {
		logger.Fatalln("error unmarshaling reverse index:", err)
	}

	var args = flag.Args()
	if len(args) == 0 {
		logger.Fatalln("At least one arg (search term) required.")
	}

	idSet := rindex[args[0]]
	for _, arg := range args[1:] {
		nextSet, ok := rindex[arg]

		// Ignore misses
		if !ok {
			continue
		}

		if *strict {
			idSet.Intersection(nextSet)
		} else {
			idSet.Union(nextSet)
		}
	}

	fmt.Println()
	for id := range idSet {
		fmt.Println("#", id)
		fmt.Println("======================")
		fmt.Printf("URL:\n%s\n\n", index[id]["url"])
		fmt.Printf("Transcript:\n%v\n\n", index[id]["transcript"])
	}
}
