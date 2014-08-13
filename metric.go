package hmetrics2

type Metric interface {
	StatAndClear() map[string]float64
}
