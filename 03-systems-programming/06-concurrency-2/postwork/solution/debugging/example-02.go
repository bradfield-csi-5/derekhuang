package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	var urls = []string{
		"https://bradfieldcs.com/courses/architecture/",
		"https://bradfieldcs.com/courses/networking/",
		"https://bradfieldcs.com/courses/databases/",
	}
	var wg sync.WaitGroup
	for i := range urls {
		// The problem is that wg.Add must happen BEFORE we launch the goroutine;
		// otherwise, it's possible that the subsequent call to wg.Wait() will be
		// waiting on nothing.
		wg.Add(1)
		go func(i int) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()

			_, err := http.Get(urls[i])
			if err != nil {
				panic(err)
			}

			fmt.Println("Successfully fetched", urls[i])
		}(i)
	}

	// Wait for all url fetches
	wg.Wait()
	fmt.Println("all url fetches done!")
}