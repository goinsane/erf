package erf

import (
	"runtime"
)

// PC returns program counters by using runtime.Callers.
func PC(size, skip int) []uintptr {
	pc := make([]uintptr, size)
	pc = pc[:runtime.Callers(skip, pc)]
	return pc
}
