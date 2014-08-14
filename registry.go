package hmetrics2

import (
	"sync"
	"time"
)

type registry struct {
	metrics map[string]Metric
	period  time.Duration
	hooks   []func(map[string]float64)
	mu      sync.Mutex
}

var hRegistry registry

func init() {
	hRegistry.init()
}

func (self *registry) init() {
	self.metrics = make(map[string]Metric)
	self.period = time.Minute
	go self.ticker()
}

func (self *registry) ticker() {
	for {
		t0 := time.Now()

		self.mu.Lock()
		period := self.period
		self.clear()
		self.mu.Unlock()

		Δt := time.Since(t0)
		if period > Δt {
			time.Sleep(period - Δt)
		}
	}
}

func (self *registry) clear() {
	data := make(map[string]float64)
	for key, metric := range self.metrics {
		metricData := metric.StatAndClear()
		for subKey, val := range metricData {
			data[key+"."+subKey] = val
		}
	}

	for _, hook := range self.hooks {
		hook(data)
	}
}

func SetPeriod(period time.Duration) {
	hRegistry.mu.Lock()
	defer hRegistry.mu.Unlock()
	hRegistry.period = period
}

func AddHook(hook func(map[string]float64)) {
	hRegistry.mu.Lock()
	defer hRegistry.mu.Unlock()
	hRegistry.hooks = append(hRegistry.hooks, hook)
}

// Register global metric and returns it and error
func RegisterGlobalMetric(name string, metric Metric) (Metric, error) {
	hRegistry.mu.Lock()
	defer hRegistry.mu.Unlock()
	if _, found := hRegistry.metrics[name]; found {
		return metric, newMetricAlreadyExistsError(name)
	}
	hRegistry.metrics[name] = metric
	return metric, nil
}

// Register global metric and returns it or panic
func MustRegisterGlobalMetric(name string, metric Metric) Metric {
	if _, err := RegisterGlobalMetric(name, metric); err != nil {
		panic(err.Error())
	}
	return metric
}

func RegisterPackageMetric(name string, metric Metric) (Metric, error) {
	pkgName := getCallerPackage()
	var metricKey string
	if pkgName != "" {
		metricKey = pkgName + "." + name
	} else {
		metricKey = name
	}
	return RegisterGlobalMetric(metricKey, metric)
}

func MustRegisterPackageMetric(name string, metric Metric) Metric {
	if _, err := RegisterPackageMetric(name, metric); err != nil {
		panic(err.Error())
	}
	return metric
}
