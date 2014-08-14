HMETRICS2
=========

Build Status
------------

![Build Status](https://travis-ci.org/orofarne/hmetrics2.svg)


Install
-------
```bash
go get 'github.com/orofarne/hmetrics2'
```

API
---
https://gowalker.org/github.com/orofarne/hmetrics2

Example
-------

    package main

    import (
            "expvar"
            "log"
            "time"

            . "github.com/orofarne/hmetrics2"
            "github.com/orofarne/hmetrics2/expvarexport"
    )

    func main() {
            HRegistry.SetPeriod(10 * time.Second)
            HRegistry.AddHook(expvarexport.Exporter("test"))
            h := NewHistogram()
            HRegistry.MustRegisterPackageMetric("my_metric", h)
            for {
                    h.AddPoint(3.14)
                    log.Print(expvar.Get("test"))
                    time.Sleep(time.Second)
            }
    }

