package myvar

import (
	"utils/base/util"
)

func (v *Var) ListItemValues(key interface{}) (values []interface{}) {
	return util.ListItemValues(v.Val(), key)
}

func (v *Var) ListItemValuesUnique(key string) []interface{} {
	return util.ListItemValuesUnique(v.Val(), key)
}
