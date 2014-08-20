package hmetrics2

import (
	"sync"
)

type Counter struct {
	mu    sync.Mutex
	count uint64
}

// Create new counter
func NewCounter() *Counter {
	return &Counter{}
}

// Increment counter
func (self *Counter) Inc() {
	self.mu.Lock()
	self.count++
	self.mu.Unlock()
}

func (self *Counter) StatAndClear() (stat map[string]float64) {
	self.mu.Lock()
	defer self.mu.Unlock()

	stat = make(map[string]float64)
	stat["count"] = float64(self.count)
	self.count = 0

	return
}
