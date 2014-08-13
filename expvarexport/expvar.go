package expvarexport

import (
	"expvar"
	"sync"
)

func Exporter(namespace string) func(map[string]float64) {
	var mu sync.Mutex
	var data map[string]float64

	expvar.Publish(namespace, expvar.Func(func() interface{} {
		mu.Lock()
		defer mu.Unlock()
		return data
	}))

	return func(newData map[string]float64) {
		mu.Lock()
		data = newData
		mu.Unlock()
	}
}
