package verror

import "runtime"

type stack []uintptr

const (
	gMAX_STACK_DEPTH = 32
)

func callers(skip ...int) stack {
	var (
		pcs [gMAX_STACK_DEPTH]uintptr
		n   = 3
	)
	if len(skip) > 0 {
		n += skip[0]
	}
	return pcs[:runtime.Callers(n, pcs[:])]
}
