package hmetrics2

import (
	. "gopkg.in/check.v1"
)

type HistogramSuite struct{}

var _ = Suite(&HistogramSuite{})

func (s *HistogramSuite) TestHistogramOn10Numbers(c *C) {
	var hist = NewHistogram()
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}

	for _, point := range data {
		hist.AddPoint(point)
	}

	var stat map[string]float64 = hist.StatAndClear()
	c.Check(len(stat), Equals, 11)
	c.Check(stat["count"], Near, 10.0, 0.0001)
	c.Check(stat["avg"], Near, ((1.0 + 10.0) / 2.0), 0.0001)
	c.Check(stat["max"], Near, 10.0, 0.0001)
	c.Check(stat["min"], Near, 1.0, 0.0001)
	// R:
	// > quantile(c(1, 2, 3, 4, 5, 6, 7, 8, 9, 10), c(0.5, 0.75, 0.95, 0.99, 0.999, 1))
	// 50%    75%    95%    99%    99.9%   100%
	// 5.500  7.750  9.550  9.910  9.991   10.000
	c.Check(stat["percentile.0_5"], In, 5.0, 6.0)
	c.Check(stat["percentile.0_75"], In, 7.0, 8.0)
	c.Check(stat["percentile.0_95"], In, 9.0, 10.0)
	c.Check(stat["percentile.0_99"], In, 9.0, 10.0)
	c.Check(stat["percentile.0_999"], In, 9.0, 10.0)
	c.Check(stat["percentile.1"], Near, 10.0, 0.0001)
	//
	c.Check(stat["rps"] > 10.0, Equals, true)
}
