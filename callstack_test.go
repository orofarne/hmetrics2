package hmetrics2

import (
	. "gopkg.in/check.v1"
)

type CallstackSuite struct{}

var _ = Suite(&CallstackSuite{})

func (s *CallstackSuite) TestCallstackNaive(c *C) {
	var pkgName string = func() string {
		return getCallerPackage()
	}()
	c.Check(pkgName, Matches, ".*hmetrics2")
}
