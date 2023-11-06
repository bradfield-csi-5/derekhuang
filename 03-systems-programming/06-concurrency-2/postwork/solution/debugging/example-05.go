package main

import (
	"fmt"
	"sync"
)

type coordinator struct {
	lock   sync.RWMutex
	leader string
}

func newCoordinator(leader string) *coordinator {
	return &coordinator{
		lock:   sync.RWMutex{},
		leader: leader,
	}
}

// c.lock must be held (in either mode) before calling logStateInternal.
func (c *coordinator) logStateInternal() {
	fmt.Printf("leader = %q\n", c.leader)
}

func (c *coordinator) logState() {
	c.lock.RLock()
	defer c.lock.RUnlock()

	c.logStateInternal()
}

func (c *coordinator) setLeader(leader string, shouldLog bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.leader = leader

	if shouldLog {
		// The problem is that sync.RWMutex is NOT re-entrant; if a goroutine
		// holds `lock` in exclusive mode but also tries to acquire it in
		// shared mode, it will block forever. The solution is to separate
		// the logging functionality into a part that doesn't involve locking,
		// and only use that part if `lock` is already held.
		c.logStateInternal()
	}
}

func main() {
	c := newCoordinator("us-east")
	c.logState()
	c.setLeader("us-west", true)
}
