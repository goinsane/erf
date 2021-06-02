package erf

import (
	"go/build"
	"strings"
)

func trimSrcPath(s string) string {
	var r string
	r = strings.TrimPrefix(s, build.Default.GOROOT+"/src/")
	if r != s {
		return r
	}
	r = strings.TrimPrefix(s, build.Default.GOPATH+"/src/")
	if r != s {
		return r
	}
	return s
}

func trimDirs(s string) string {
	for i := len(s) - 1; i > 0; i-- {
		if s[i] == '/' {
			return s[i+1:]
		}
	}
	return s
}
