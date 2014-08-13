package hmetrics2

import (
	"fmt"
	"math"
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

// Near checker.

type nearChecker struct {
	*CheckerInfo
}

// The Near checker verifies that the obtained value is near to
// the expected value with ε.
//
// For example:
//
//     c.Assert(value, Near, 42.0, 0.001)
//
var Near Checker = &nearChecker{
	&CheckerInfo{Name: "Near", Params: []string{"obtained", "expected", "ε"}},
}

func (checker *nearChecker) Check(params []interface{}, names []string) (result bool, error string) {
	defer func() {
		if v := recover(); v != nil {
			result = false
			error = fmt.Sprint(v)
		}
	}()
	ε := params[2].(float64)
	result = math.Abs(params[0].(float64)-params[1].(float64)) < ε
	return
}

// In checker.

type inChecker struct {
	*CheckerInfo
}

// The In checker verifies that the obtained value in interval.
//
// For example:
//
//     c.Assert(value, In, 42.0, 57.0)
//
var In Checker = &inChecker{
	&CheckerInfo{Name: "In", Params: []string{"obtained", "min", "max"}},
}

func (checker *inChecker) Check(params []interface{}, names []string) (result bool, error string) {
	defer func() {
		if v := recover(); v != nil {
			result = false
			error = fmt.Sprint(v)
		}
	}()
	result = params[0].(float64) >= params[1].(float64) && params[0].(float64) <= params[2].(float64)
	return
}
