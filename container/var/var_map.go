package vvar

import "utils/os/conv"

func (v *Var) Map(tags ...string) map[string]interface{} {
	return conv.Map(v.Val(), tags...)
}

func (v *Var) MapStrAny() map[string]interface{} {
	return v.Map()
}

func (v *Var) MapStrStr(tags ...string) map[string]string {
	return conv.MapStrStr(v.Val(), tags...)
}

func (v *Var) MapStrVar(tags ...string) map[string]*Var {
	m := v.Map(tags...)
	if len(m) > 0 {
		vMap := make(map[string]*Var, len(m))
		for k, v := range m {
			vMap[k] = New(v)
		}
		return vMap
	}
	return nil
}

func (v *Var) MapDeep(tags ...string) map[string]interface{} {
	return conv.MapDeep(v.Val(), tags...)
}

func (v *Var) MapStrStrDeep(tags ...string) map[string]string {
	return conv.MapStrStrDeep(v.Val(), tags...)
}

func (v *Var) MapStrVarDeep(tags ...string) map[string]*Var {
	m := v.MapDeep(tags...)
	if len(m) > 0 {
		vMap := make(map[string]*Var, len(m))
		for k, v := range m {
			vMap[k] = New(v)
		}
		return vMap
	}
	return nil
}

func (v *Var) Maps(tags ...string) []map[string]interface{} {
	return conv.Maps(v.Val(), tags...)
}

func (v *Var) MapToMap(pointer interface{}, mapping ...map[string]string) (err error) {
	return conv.MapToMap(v.Val(), pointer, mapping...)
}

func (v *Var) MapToMapDeep(pointer interface{}, mapping ...map[string]string) (err error) {
	return conv.MapToMapDeep(v.Val(), pointer, mapping...)
}

func (v *Var) MapToMaps(pointer interface{}, mapping ...map[string]string) (err error) {
	return conv.MapToMaps(v.Val(), pointer, mapping...)
}

func (v *Var) MapToMapsDeep(pointer interface{}, mapping ...map[string]string) (err error) {
	return conv.MapToMapsDeep(v.Val(), pointer, mapping...)
}
