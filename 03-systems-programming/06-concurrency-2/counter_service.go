// Author: Patch Neranartkomol

package counterservice

import (
	"sync"
	"sync/atomic"
)

type CounterService interface {
	// Returns values in ascending order; it should be safe to call
	// getNext() concurrently from multiple goroutines without any
	// additional synchronization on the caller's side.
	getNext() uint64
}

type UnsynchronizedCounterService struct {
	count uint64
}

// getNext() - This one can be UNSAFE
func (counter *UnsynchronizedCounterService) getNext() uint64 {
	counter.count += 1
	return counter.count
}

type AtomicCounterService struct {
	count atomic.Uint64
}

// getNext() with sync/atomic
func (counter *AtomicCounterService) getNext() uint64 {
	return counter.count.Add(1)
}

type MutexCounterService struct {
	count uint64
	m     sync.Mutex
}

// getNext() with sync/Mutex
func (counter *MutexCounterService) getNext() uint64 {
	counter.m.Lock()
	counter.count += 1
	counter.m.Unlock()
	return counter.count
}

type ChannelCounterService struct {
	count uint64
}

// A constructor for ChannelCounterService
func newChannelCounterService() *ChannelCounterService {
	cs := ChannelCounterService{}
	return &cs
}

// getNext() with goroutines and channels
func (counter *ChannelCounterService) getNext() uint64 {
	panic("getNext not implemented")
}
