package erf

import (
	"go/build"
	"os"
	"runtime"
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
		if s[i] == '/' || s[i] == os.PathSeparator {
			return s[i+1:]
		}
	}
	return s
}

func getPC(size, skip int) []uintptr {
	pc := make([]uintptr, size)
	pc = pc[:runtime.Callers(skip, pc)]
	return pc
}
