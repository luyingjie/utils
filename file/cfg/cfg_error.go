package cfg

import (
	"github.com/luyingjie/utils/util/cmdenv"
)

const (
	// ERROR_PRINT_KEY is used to specify the key controlling error printing to stdout.
	// This error is designed not to be returned by functions.
	ERROR_PRINT_KEY = "cfg.errorprint"
)

// errorPrint checks whether printing error to stdout.
func errorPrint() bool {
	return cmdenv.Get(ERROR_PRINT_KEY, true).Bool()
}
