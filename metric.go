package hmetrics2

type Metric interface {
	Stat() map[string]float64
}
