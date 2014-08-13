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

var HRegistry registry

func init() {
	HRegistry.init()
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

func (self *registry) SetPeriod(period time.Duration) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.period = period
}

func (self *registry) AddHook(hook func(map[string]float64)) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.hooks = append(self.hooks, hook)
}

func (self *registry) RegisterGlobalMetric(name string, metric Metric) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	if _, found := self.metrics[name]; found {
		return newMetricAlreadyExistsError(name)
	}
	self.metrics[name] = metric
	return nil
}

func (self *registry) MustRegisterGlobalMetric(name string, metric Metric) {
	if err := self.RegisterGlobalMetric(name, metric); err != nil {
		panic(err.Error())
	}
}

func (self *registry) RegisterPackageMetric(name string, metric Metric) error {
	pkgName := getCallerPackage()
	var metricKey string
	if pkgName != "" {
		metricKey = "pkgName" + "." + name
	} else {
		metricKey = name
	}
	return self.RegisterGlobalMetric(metricKey, metric)
}

func (self *registry) MustRegisterPackageMetric(name string, metric Metric) {
	if err := self.RegisterPackageMetric(name, metric); err != nil {
		panic(err.Error())
	}
}
