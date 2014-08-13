package hmetrics2

import (
	. "gopkg.in/check.v1"
)

type ErrSuite struct{}

var _ = Suite(&ErrSuite{})

func (s *ErrSuite) TestErrorMessage(c *C) {
	var err error = newMetricAlreadyExistsError("metric_name")

	c.Check(err.Error(), Equals, "Metric \"metric_name\" already exists")
}
