/*
Exercise 1.7: The function call io.Copy(dst, src) reads from src and writes to dst. Use it instead of ioutil.ReadAll
to copy the response body to os.Stdout without requiring a buffer large enough to hold the entire stream. Be sure to
check the error result of io.Copy.

Exercise 1.8: Modify fetch to add the prefix http:// to each argument URL if it is missing.
You might want to use strings.HasPrefix.

Exercise 1.9: Modify fetch to also print the HTTP status code, found in resp.Status.
*/

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	for _, url := range os.Args[1:] {
		if !strings.HasPrefix(url, "http://") {
			url = "http://" + url
		}
		fmt.Printf("Getting url: %s\n", url)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Received status code: %s\n", resp.Status)
		if _, err = io.Copy(os.Stdout, resp.Body); err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}
		resp.Body.Close()
	}
}
