package main

import (
	"expvar"
	"log"
	"time"

	hmetrics2 ".."
	"../expvarexport"
)

func main() {
	hmetrics2.SetPeriod(10 * time.Second)
	hmetrics2.AddHook(expvarexport.Exporter("test"))
	h := hmetrics2.NewHistogram()
	hmetrics2.MustRegisterPackageMetric("my_metric", h)
	for {
		h.AddPoint(3.14)
		log.Print(expvar.Get("test"))
		time.Sleep(time.Second)
	}
}
