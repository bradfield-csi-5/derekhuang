package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

const (
	fileName = "xkcd.json"
	max      = 2800
	xkcdUrl  = "https://xkcd.com"
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

type Record map[string]string

type item struct {
	comic Comic
	err   error
}

func PopulateIndex(logger *log.Logger) {
	var index Index
	var skipped int = 0
	var wg sync.WaitGroup

	file, err := os.OpenFile(fileName, os.O_APPEND, 0666)
	if err != nil {
		logger.Printf("Failed to open file. Creating a new one...\n")
		file, err = os.Create(fileName)
		if err != nil {
			logger.Fatalln("failed to create new file:", err)
		}
	}

	if err := json.NewDecoder(file).Decode(&index); err != nil {
		if err != io.EOF {
			logger.Fatalln("error decoding json:", err)
		}
		// The index was malformed/nonexistent/empty so create a new one
		index = make(Index)
	}

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

	if skipped < max {
		go func() {
			wg.Wait()
			close(ch)
		}()

		logger.Println("Building index...")
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

		logger.Println("Encoding json...")
		b, err := json.MarshalIndent(index, "", "  ")
		if err != nil {
			logger.Fatalln("error encoding json:", err)
		}

		logger.Println("Writing to disk...")
		if err := os.WriteFile(fileName, b, 0666); err != nil {
			logger.Fatalln("error writing json:", err)
		}
	}

	logger.Println("Done")
}

func fetch(i int, logger *log.Logger) (comic Comic, err error) {
	logger.Printf("Fetching #%d...\n", i)
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

	if result.Num == 0 {
		result.Num = i
	}

	return result, nil
}
