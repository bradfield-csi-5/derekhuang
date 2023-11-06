package main

import (
	"fmt"
)

const numTasks = 3

func main() {
	// The problem is that sending or receiving on a nil channel always blocks!
	// We need to initialize the channel.
	done := make(chan struct{})
	for i := 0; i < numTasks; i++ {
		go func() {
			fmt.Println("running task...")

			// Signal that task is done
			done <- struct{}{}
		}()
	}

	// Wait for tasks to complete
	for i := 0; i < numTasks; i++ {
		<-done
	}
	fmt.Printf("all %d tasks done!\n", numTasks)
}
