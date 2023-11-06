package main

import (
	"fmt"
	"sync"
	"time"
)

type dbService struct {
	lock       *sync.RWMutex
	connection string
}

func newDbService(connection string) *dbService {
	return &dbService{
		lock:       &sync.RWMutex{},
		connection: connection,
	}
}

// Initially, you might think that the re-entrancy issue from
// example-05.go doesn't apply here, because takeSnapshot() and logState()
// both use `lock` in shared mode. However, it's still a problem because of
// this scenario:
//   - takeSnapshot acquires `lock` in shared mode
//   - updateConnection tries to acquire `lock` in exclusive mode; crucially,
//     this forces future attempts to acquire `lock` in shared mode to block
//   - takeSnapshot calls into `logState`, which tries to acquire `lock` in
//     shared mode; it cannot, because there is already another goroutine waiting
//     to acquire it in exclusive mode
func (d *dbService) logStateInternal() {
	fmt.Printf("connection %q is healthy\n", d.connection)
}

func (d *dbService) logState() {
	d.lock.RLock()
	defer d.lock.RUnlock()

	d.logStateInternal()
}

func (d *dbService) takeSnapshot() {
	d.lock.RLock()
	defer d.lock.RUnlock()

	fmt.Printf("Taking snapshot over connection %q\n", d.connection)

	// Simulate slow operation
	time.Sleep(time.Second)

	d.logStateInternal()
}

func (d *dbService) updateConnection(connection string) {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.connection = connection
}

func main() {
	d := newDbService("127.0.0.1:3001")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		d.takeSnapshot()
	}()

	// Simulate other DB accesses
	time.Sleep(200 * time.Millisecond)

	wg.Add(1)
	go func() {
		defer wg.Done()

		d.updateConnection("127.0.0.1:8080")
	}()

	wg.Wait()
}
