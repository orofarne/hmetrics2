package hmetrics2

import (
	"runtime"
	"strings"
)

func getCallerPackage() string {
	pc, _, _, ok := runtime.Caller(3)
	if !ok {
		return ""
	}
	f := runtime.FuncForPC(pc)
	fName := f.Name()
	packageName := fName[:strings.Index(fName, ".")]
	return packageName
}
