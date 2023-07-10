package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	fileName = "xkcd.json"
	max      = 2799 // TODO: change to 2799
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

func PopulateIndex(logger *log.Logger) {
	var record map[int]map[string]string

	file, err := os.OpenFile(fileName, os.O_APPEND, 0666)
	if err != nil {
		logger.Printf("Failed to open file. Creating a new one...\n")
		file, err = os.Create(fileName)
		if err != nil {
			logger.Fatalln("failed to create new file:", err)
		}
	}

	if err := json.NewDecoder(file).Decode(&record); err != nil {
		if err != io.EOF {
			logger.Fatalln("error decoding json:", err)
		}
		// The index was malformed/nonexistent/empty so create a new one
		record = make(map[int]map[string]string)
	}

	for i := 1; i <= max; i++ {
		if _, ok := record[i]; ok {
			logger.Printf("Comic #%d found in index. Skipping...\n", i)
			continue
		}

		resp, err := http.Get(fmt.Sprintf("%s/%d/%s", xkcdUrl, i, "info.0.json"))
		if err != nil {
			logger.Fatalln("get error:", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			logger.Printf("get request failed for comic #%d: %s\n", i, resp.Status)
			continue
		}

		var result Comic
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			logger.Fatalln("json decode failed:", err)
		}

		record[i] = map[string]string{
			"url":        fmt.Sprintf("%s/%d", xkcdUrl, i),
			"transcript": fmt.Sprintf("%#v", result.Transcript),
		}
	}

	b, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		logger.Fatalln("error encoding json:", err)
	}

	if err := os.WriteFile(fileName, b, 0666); err != nil {
		logger.Fatalln("error writing json:", err)
	}
}
