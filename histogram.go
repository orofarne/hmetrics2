package hmetrics2

import (
	"fmt"
	"math"
	"sort"
	"sync"
)

type Histogram struct {
	mu     sync.Mutex
	values []float64
	min    float64
	max    float64
	sum    float64
	count  uint64
}

// Create new metric
func NewHistogram() *Histogram {
	return &Histogram{
		min: math.Inf(1),
		max: math.Inf(-1),
	}
}

// Add point to histogram
// Not threadsafe!
func (self *Histogram) AddPoint(val float64) {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.values = append(self.values, val)
	if val < self.min {
		self.min = val
	}
	if val > self.max {
		self.max = val
	}
	self.sum += val
	self.count++
}

func (self *Histogram) clear() {
	self.values = nil
	self.min = math.Inf(1)
	self.max = math.Inf(-1)
	self.sum = 0.0
	self.count = 0
}

func (self *Histogram) percentiles(ps []float64) []float64 {
	scores := make([]float64, len(ps))
	size := len(self.values)
	if size > 0 {
		sort.Float64s(self.values)
		for i, p := range ps {
			pos := p * float64(size)
			if pos < 1.0 {
				scores[i] = float64(self.values[0])
			} else if pos >= float64(size) {
				scores[i] = float64(self.values[size-1])
			} else {
				lower := float64(self.values[int(pos)-1])
				upper := float64(self.values[int(pos)])
				scores[i] = lower + (pos-math.Floor(pos))*(upper-lower)
			}
		}
	}
	return scores
}

func (self *Histogram) StatAndClear() (stat map[string]float64) {
	self.mu.Lock()
	defer self.mu.Unlock()

	stat = make(map[string]float64)
	// Basic statistics
	stat["min"] = self.min
	stat["max"] = self.max
	stat["avg"] = self.sum / float64(self.count)
	stat["count"] = float64(self.count)
	// Percentiles
	percs := []float64{0.5, 0.75, 0.95, 0.99, 0.999, 1.0}
	percsValues := self.percentiles(percs)
	for i, p := range percsValues {
		stat[fmt.Sprintf("percentile_%v", percs[i])] = p
	}
	// Clear data
	self.clear()
	return
}
