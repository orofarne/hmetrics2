package hmetrics2

import (
	"fmt"
)

type MetricAlreadyExists struct {
	metricName string
}

func newMetricAlreadyExistsError(name string) *MetricAlreadyExists {
	return &MetricAlreadyExists{
		metricName: name,
	}
}

func (self *MetricAlreadyExists) Error() string {
	return fmt.Sprintf("Metric %s already exists", self.metricName)
}
