package expvarexport

import (
	"expvar"
	"math"
	"sync"
)

func Exporter(namespace string) func(map[string]float64) {
	var mu sync.Mutex
	var data = make(map[string]*float64)

	expvar.Publish(namespace, expvar.Func(func() interface{} {
		mu.Lock()
		defer mu.Unlock()
		return data
	}))

	return func(newData map[string]float64) {
		mu.Lock()
		data = make(map[string]*float64)
		for k, v := range newData {
			if math.IsNaN(v) || math.IsInf(v, 0) {
				data[k] = nil
			} else {
				val := new(float64)
				*val = v
				data[k] = val
			}
		}
		mu.Unlock()
	}
}
