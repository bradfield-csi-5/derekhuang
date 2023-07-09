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
	max      = 1 // TODO: change to 2799
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

func PopulateIndex() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var record map[int]map[string]string

	file, err := os.OpenFile(fileName, os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("error opening file: ", err)
	}
	if err := json.NewDecoder(file).Decode(&record); err != nil {
		if err != io.EOF {
			log.Fatalln("error decoding json: ", err)
		}
		// the index didn't exist/didn't have any data
		record = make(map[int]map[string]string)
	}

	for i := 1; i <= max; i++ {
		if _, ok := record[i]; ok {
			log.Printf("Comic #%d found in index. Skipping...\n", i)
			continue
		}
		resp, err := http.Get(fmt.Sprintf("%s/%d/%s", xkcdUrl, i, "info.0.json"))
		if err != nil {
			log.Fatalln("get error: ", err)
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			log.Printf("get request failed for comic #%d: %s\n", i, resp.Status)
			continue
		}
		var result Comic
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			log.Fatalln("json decode failed: ", err)
		}
		record[i] = map[string]string{
			"url":        fmt.Sprintf("%s/%d", xkcdUrl, i),
			"transcript": fmt.Sprintf("%#v", result.Transcript),
		}
	}

	b, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		log.Fatalln("error encoding json: ", err)
	}
	if err := os.WriteFile(fileName, b, 0666); err != nil {
		log.Fatalln("error writing json: ", err)
	}
}
