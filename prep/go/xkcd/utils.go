package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

type Comic struct {
	Month      string
	Num        int
	Link       string
	Year       string
	News       string
	SafeTitle  string
	Transcript string
	Alt        string
	Img        string
	Title      string
	Day        string
}

type Index map[int]Record

type ReverseIndex map[string]IntSet

type Record map[string]string

type IntSet map[int]bool

const (
	IndexFileName    = "xkcd.json"
	RevIndexFileName = "dckx.json"
	max              = 2800
	xkcdUrl          = "https://xkcd.com"
)

var ignoreRe = regexp.MustCompile("Title text:")

type item struct {
	comic Comic
	err   error
}

func PopulateIndex(logger *log.Logger) {
	var index Index
	var skipped int = 0
	var wg sync.WaitGroup

	file, err := os.OpenFile(IndexFileName, os.O_APPEND, 0666)
	if err != nil {
		logger.Printf("Failed to read index. Creating a new one...\n")
		file, err = os.Create(IndexFileName)
		if err != nil {
			logger.Fatalln("error creating new file:", err)
		}
	}

	if err := json.NewDecoder(file).Decode(&index); err != nil {
		if err != io.EOF {
			logger.Fatalln("error decoding json:", err)
		}
		// The index was malformed/nonexistent/empty so create a new one
		index = make(Index)
	}

	// Short circuit if the index is fully populated
	// (compare with 1 less than max to account for #404)
	if len(index) == max-1 {
		return
	}

	logger.Println("Fetching missing comics...")
	ch := make(chan item, max)
	for i := 1; i <= max; i++ {
		if i == 404 { // 404 Not Found
			skipped++
			continue
		}
		if _, ok := index[i]; ok {
			skipped++
			continue
		}

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var it item
			it.comic, it.err = fetch(i, logger)
			if it.err != nil {
				logger.Printf("error getting #%d\n", it.comic.Num)
				return
			}
			ch <- it
		}(i)
	}
	logger.Println("  done")

	if skipped < max {
		go func() {
			wg.Wait()
			close(ch)
		}()

		logger.Println("Building index in memory...")
		for it := range ch {
			index[it.comic.Num] = Record{
				"url":        fmt.Sprintf("%s/%d", xkcdUrl, it.comic.Num),
				"transcript": it.comic.Transcript,
				"alt":        it.comic.Alt,
				"day":        it.comic.Day,
				"month":      it.comic.Month,
				"year":       it.comic.Year,
				"num":        strconv.Itoa(it.comic.Num),
				"link":       it.comic.Link,
				"img":        it.comic.Img,
				"news":       it.comic.News,
				"title":      it.comic.Title,
				"safe_title": it.comic.SafeTitle,
			}
		}

		b, err := json.MarshalIndent(index, "", "  ")
		if err != nil {
			logger.Fatalln("error encoding json:", err)
		}

		logger.Println("Writing index to disk...")
		if err := os.WriteFile(IndexFileName, b, 0666); err != nil {
			logger.Fatalln("error writing json:", err)
		}

		logger.Println("  done")
	}
}

func fetch(i int, logger *log.Logger) (comic Comic, err error) {
	resp, err := http.Get(fmt.Sprintf("%s/%d/info.0.json", xkcdUrl, i))
	if err != nil {
		return Comic{Num: i}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Printf("get request failed for #%d: %s\n", i, resp.Status)
		return Comic{Num: i}, errors.New(resp.Status)
	}

	var result Comic
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Comic{Num: i}, err
	}

	// Fill in num in case it was empty in the response
	if result.Num == 0 {
		result.Num = i
	}

	return result, nil
}

func BuildReverseIndex(logger *log.Logger) {
	if _, err := os.Stat(RevIndexFileName); err == nil {
		return
	}

	b, err := os.ReadFile(IndexFileName)
	if err != nil {
		logger.Fatalln("error reading index:", err)
	}

	var index Index
	err = json.Unmarshal(b, &index)
	if err != nil {
		logger.Fatalln("error unmarshaling index:", err)
	}

	logger.Printf("Normalizing words from index...")
	rindex := make(ReverseIndex)
	for i, cmc := range index {
		normalized := normalize(cmc["alt"], cmc["safe_title"], cmc["transcript"])
		for _, word := range normalized {
			if _, ok := rindex[word]; !ok {
				rindex[word] = make(IntSet)
			}
			rindex[word][i] = true
		}
	}

	b, err = json.MarshalIndent(rindex, "", "  ")
	if err != nil {
		logger.Fatalln("error marshaling reverse index:", err)
	}

	logger.Printf("Writing reverse index to disk...")
	if err := os.WriteFile(RevIndexFileName, b, 0666); err != nil {
		logger.Fatalln("error writing reverse index:", err)
	}

	logger.Println("  done")
}

func normalize(strs ...string) []string {
	var ret []string

	for _, s := range strs {
		// Remove ignored words and convert to lowercase
		normalizedStr := ignoreRe.ReplaceAllString(strings.ToLower(s), "")

		// Remove everything that isn't an ascii letter
		normalizedWords := strings.FieldsFunc(normalizedStr, func(c rune) bool {
			return !unicode.IsLetter(c) || (c > unicode.MaxASCII)
		})

		// Ignore short words
		for _, word := range normalizedWords {
			if len(word) > 2 {
				ret = append(ret, word)
			}
		}
	}

	return ret
}

func Intersection(s1, s2 IntSet) IntSet {
	ret := make(IntSet)
	for k := range s1 {
		if _, exists := s2[k]; exists {
			ret[k] = true
		}
	}
	return ret
}

func Union(s1, s2 IntSet) IntSet {
	for k := range s2 {
		if _, exists := s1[k]; !exists {
			s1[k] = true
		}
	}
	return s1
}
