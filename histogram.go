package hmetrics2

import (
	"fmt"
	"math"
	"sort"
)

type Histogram struct {
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

func (self *Histogram) Clear() {
	self.values = nil
	self.min = math.Inf(1)
	self.max = math.Inf(-1)
	self.sum = 0.0
	self.count = 0
}

func (self *Histogram) Min() float64 {
	return self.min
}

func (self *Histogram) Max() float64 {
	return self.max
}

func (self *Histogram) Avg() float64 {
	return self.sum / float64(self.count)
}

func (self *Histogram) Count() uint64 {
	return self.count
}

func (self *Histogram) Percentiles(ps []float64) []float64 {
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

func (self *Histogram) Stat() (stat map[string]float64) {
	stat = make(map[string]float64)
	// Basic statistics
	stat["min"] = self.Min()
	stat["max"] = self.Max()
	stat["avg"] = self.Avg()
	stat["count"] = float64(self.Count())
	// Percentiles
	percs := []float64{0.5, 0.75, 0.95, 0.99, 0.999, 1.0}
	percsValues := self.Percentiles(percs)
	for i, p := range percsValues {
		stat[fmt.Sprintf("percentile_%v", percs[i])] = p
	}
	return
}
