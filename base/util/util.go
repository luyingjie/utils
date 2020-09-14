package util

import (
	"utils/utils/empty"
)

func Throw(exception interface{}) {
	panic(exception)
}

func TryCatch(try func(), catch ...func(exception interface{})) {
	defer func() {
		if e := recover(); e != nil && len(catch) > 0 {
			catch[0](e)
		}
	}()
	try()
}

func IsEmpty(value interface{}) bool {
	return empty.IsEmpty(value)
}
