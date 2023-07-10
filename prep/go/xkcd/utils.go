package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	fileName = "xkcd.json"
	max      = 2799
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

func PopulateIndex(logger *log.Logger) {
	var index Index
	var skipped uint = 0

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

	for i := 1; i <= max; i++ {
		if i == 404 { // 404 Not Found
			skipped++
			continue
		}
		if _, ok := index[i]; ok {
			logger.Printf("Comic #%d found in index. Skipping...\n", i)
			skipped++
			continue
		}

		logger.Printf("Fetching #%d...\n", i)
		resp, err := http.Get(fmt.Sprintf("%s/%d/%s", xkcdUrl, i, "info.0.json"))
		if err != nil {
			logger.Fatalf("get error for #%d: %v\n", i, err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			logger.Printf("get request failed for #%d: %s\n", i, resp.Status)
			continue
		}

		var result Comic
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			logger.Fatalf("json decode failed for #%d: %v\n", i, err)
		}

		index[i] = Record{
			"url":        fmt.Sprintf("%s/%d", xkcdUrl, i),
			"transcript": fmt.Sprintf("%#v", result.Transcript),
			"alt":        fmt.Sprintf("%#v", result.Alt),
			"day":        result.Day,
			"month":      result.Month,
			"year":       result.Year,
			"num":        strconv.Itoa(result.Num),
			"link":       result.Link,
			"img":        result.Img,
			"news":       result.News,
			"title":      fmt.Sprintf("%#v", result.Title),
			"safe_title": result.SafeTitle,
		}
	}

	if skipped < max {
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
