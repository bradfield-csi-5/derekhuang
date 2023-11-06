package main

import (
	"fmt"
	"math/rand"
	"time"
)

var responses = []string{
	"200 OK",
	"402 Payment Required",
	"418 I'm a teapot",
}

func randomDelay(maxMillis int) time.Duration {
	return time.Duration(rand.Intn(maxMillis)) * time.Millisecond
}

func query(endpoint string) string {
	// Simulate querying the given endpoint
	delay := randomDelay(100)
	time.Sleep(delay)

	i := rand.Intn(len(responses))
	return responses[i]
}

// Query each of the mirrors in parallel and return the first
// response (this approach increases the amount of traffic but
// significantly improves "tail latency")
func parallelQuery(endpoints []string) string {
	// The problem is that `parallelQuery` will only read one value
	// out of the results channel, but there are two other goroutines
	// trying to send on this channel. The result is a "goroutine leak",
	// where the server process accumulates more and more goroutines
	// that will be blocked forever. We can fix this by changing `results`
	// into a buffered channel with a large enough size.
	results := make(chan string, len(endpoints))
	for i := range endpoints {
		go func(i int) {
			results <- query(endpoints[i])
		}(i)
	}
	return <-results
}

func main() {
	var endpoints = []string{
		"https://fakeurl.com/endpoint",
		"https://mirror1.com/endpoint",
		"https://mirror2.com/endpoint",
	}

	// Simulate long-running server process.
	// Hint: What will happen to the server's memory usage if it runs
	// continuously for a very long time?
	for {
		fmt.Println(parallelQuery(endpoints))
		delay := randomDelay(100)
		time.Sleep(delay)
	}
}
