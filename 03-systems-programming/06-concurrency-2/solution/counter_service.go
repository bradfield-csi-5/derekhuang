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
	prevMax uint64
}

func (counter *UnsynchronizedCounterService) getNext() uint64 {
	counter.prevMax++
	result := counter.prevMax
	return result
}

type AtomicCounterService struct {
	prevMax uint64
}

func (counter *AtomicCounterService) getNext() uint64 {
	return atomic.AddUint64(&counter.prevMax, 1)
}

type MutexCounterService struct {
	prevMax uint64
	mu      sync.Mutex
}

func (counter *MutexCounterService) getNext() uint64 {
	counter.mu.Lock()
	defer counter.mu.Unlock()
	counter.prevMax++
	result := counter.prevMax
	return result
}

type ChannelCounterService struct {
	prevMax uint64
	req     chan struct{}
	resp    chan uint64
}

func newChannelCounterService() *ChannelCounterService {
	cs := ChannelCounterService{
		req:  make(chan struct{}),
		resp: make(chan uint64),
	}
	// Launch monitor goroutine
	go func() {
		for range cs.req {
			cs.prevMax++
			cs.resp <- cs.prevMax
		}
	}()
	return &cs
}

func (counter *ChannelCounterService) getNext() uint64 {
	counter.req <- struct{}{}
	return <-counter.resp
}
