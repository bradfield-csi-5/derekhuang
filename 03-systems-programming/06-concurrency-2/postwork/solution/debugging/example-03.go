package main

import (
	"fmt"
)

func main() {
	done := make(chan struct{}, 1)
	go func() {
		fmt.Println("performing initialization...")
		done <- struct{}{}
	}()

	// The problem is that sending on a buffered channel won't
	// block if the buffer isn't full; we can fix the issue by
	// either making `done` an unbuffered channel, or by having
	// the main function receive instead of send
	<-done
	fmt.Println("initialization done, continuing with rest of program")
}
