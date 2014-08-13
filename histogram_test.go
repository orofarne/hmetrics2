package hmetrics2

import (
	"math"

	. "gopkg.in/check.v1"
)

type HistogramSuite struct{}

var _ = Suite(&HistogramSuite{})

func (s *HistogramSuite) TestMath(c *C) {
	c.Check(math.Inf(1) > 3.0, Equals, true)
	c.Check(math.Inf(-1) < 3.0, Equals, true)
	c.Check(math.Inf(-1) < math.Inf(1), Equals, true)
}

func (s *HistogramSuite) TestHistogramOn10Numbers(c *C) {
	var hist = NewHistogram()
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}

	for _, point := range data {
		hist.AddPoint(point)
	}

	c.Check(hist.Count(), Equals, uint64(10))
	c.Check(hist.Avg(), Near, ((1.0 + 10.0) / 2.0), 0.0001)
	c.Check(hist.Max(), Near, 10.0, 0.0001)
	c.Check(hist.Min(), Near, 1.0, 0.0001)
	c.Check(hist.Percentiles([]float64{0.5})[0], In, 5.0, 6.0)

	var stat map[string]float64 = hist.Stat()
	c.Check(len(stat), Equals, 10)
	c.Check(stat["count"], Near, 10.0, 0.0001)
	c.Check(stat["avg"], Near, ((1.0 + 10.0) / 2.0), 0.0001)
	c.Check(stat["max"], Near, 10.0, 0.0001)
	c.Check(stat["min"], Near, 1.0, 0.0001)
	// R:
	// > quantile(c(1, 2, 3, 4, 5, 6, 7, 8, 9, 10), c(0.5, 0.75, 0.95, 0.99, 0.999, 1))
	// 50%    75%    95%    99%    99.9%   100%
	// 5.500  7.750  9.550  9.910  9.991   10.000
	c.Check(stat["percentile_0.5"], In, 5.0, 6.0)
	c.Check(stat["percentile_0.75"], In, 7.0, 8.0)
	c.Check(stat["percentile_0.95"], In, 9.0, 10.0)
	c.Check(stat["percentile_0.99"], In, 9.0, 10.0)
	c.Check(stat["percentile_0.999"], In, 9.0, 10.0)
	c.Check(stat["percentile_1"], Near, 10.0, 0.0001)
}
