package json

import (
	"fmt"
	"time"

	"github.com/luyingjie/utils/util"

	"github.com/luyingjie/utils/conv"

	vtime "github.com/luyingjie/utils/os/time"

	vvar "github.com/luyingjie/utils/container/var"
)

func (j *Json) Value() interface{} {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return *(j.p)
}

func (j *Json) IsNil() bool {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.p == nil || *(j.p) == nil
}

func (j *Json) Get(pattern string, def ...interface{}) interface{} {
	j.mu.RLock()
	defer j.mu.RUnlock()

	if pattern == "" {
		return nil
	}

	if pattern == "." {
		return *j.p
	}

	var result *interface{}
	if j.vc {
		result = j.getPointerByPattern(pattern)
	} else {
		result = j.getPointerByPatternWithoutViolenceCheck(pattern)
	}
	if result != nil {
		return *result
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}

func (j *Json) GetVar(pattern string, def ...interface{}) *vvar.Var {
	return vvar.New(j.Get(pattern, def...))
}

func (j *Json) GetVars(pattern string, def ...interface{}) []*vvar.Var {
	return vvar.New(j.Get(pattern, def...)).Vars()
}

func (j *Json) GetMap(pattern string, def ...interface{}) map[string]interface{} {
	result := j.Get(pattern, def...)
	if result != nil {
		return conv.Map(result)
	}
	return nil
}

func (j *Json) GetMapStrStr(pattern string, def ...interface{}) map[string]string {
	result := j.Get(pattern, def...)
	if result != nil {
		return conv.MapStrStr(result)
	}
	return nil
}

func (j *Json) GetMaps(pattern string, def ...interface{}) []map[string]interface{} {
	result := j.Get(pattern, def...)
	if result != nil {
		return conv.Maps(result)
	}
	return nil
}

func (j *Json) GetJson(pattern string, def ...interface{}) *Json {
	return New(j.Get(pattern, def...))
}

func (j *Json) GetJsons(pattern string, def ...interface{}) []*Json {
	array := j.GetArray(pattern, def...)
	if len(array) > 0 {
		jsonSlice := make([]*Json, len(array))
		for i := 0; i < len(array); i++ {
			jsonSlice[i] = New(array[i])
		}
		return jsonSlice
	}
	return nil
}

func (j *Json) GetJsonMap(pattern string, def ...interface{}) map[string]*Json {
	m := j.GetMap(pattern, def...)
	if len(m) > 0 {
		jsonMap := make(map[string]*Json, len(m))
		for k, v := range m {
			jsonMap[k] = New(v)
		}
		return jsonMap
	}
	return nil
}

func (j *Json) GetArray(pattern string, def ...interface{}) []interface{} {
	return conv.Interfaces(j.Get(pattern, def...))
}

func (j *Json) GetString(pattern string, def ...interface{}) string {
	return conv.String(j.Get(pattern, def...))
}

func (j *Json) GetBytes(pattern string, def ...interface{}) []byte {
	return conv.Bytes(j.Get(pattern, def...))
}

func (j *Json) GetBool(pattern string, def ...interface{}) bool {
	return conv.Bool(j.Get(pattern, def...))
}

func (j *Json) GetInt(pattern string, def ...interface{}) int {
	return conv.Int(j.Get(pattern, def...))
}

func (j *Json) GetInt8(pattern string, def ...interface{}) int8 {
	return conv.Int8(j.Get(pattern, def...))
}

func (j *Json) GetInt16(pattern string, def ...interface{}) int16 {
	return conv.Int16(j.Get(pattern, def...))
}

func (j *Json) GetInt32(pattern string, def ...interface{}) int32 {
	return conv.Int32(j.Get(pattern, def...))
}

func (j *Json) GetInt64(pattern string, def ...interface{}) int64 {
	return conv.Int64(j.Get(pattern, def...))
}

func (j *Json) GetUint(pattern string, def ...interface{}) uint {
	return conv.Uint(j.Get(pattern, def...))
}

func (j *Json) GetUint8(pattern string, def ...interface{}) uint8 {
	return conv.Uint8(j.Get(pattern, def...))
}

func (j *Json) GetUint16(pattern string, def ...interface{}) uint16 {
	return conv.Uint16(j.Get(pattern, def...))
}

func (j *Json) GetUint32(pattern string, def ...interface{}) uint32 {
	return conv.Uint32(j.Get(pattern, def...))
}

func (j *Json) GetUint64(pattern string, def ...interface{}) uint64 {
	return conv.Uint64(j.Get(pattern, def...))
}

func (j *Json) GetFloat32(pattern string, def ...interface{}) float32 {
	return conv.Float32(j.Get(pattern, def...))
}

func (j *Json) GetFloat64(pattern string, def ...interface{}) float64 {
	return conv.Float64(j.Get(pattern, def...))
}

func (j *Json) GetFloats(pattern string, def ...interface{}) []float64 {
	return conv.Floats(j.Get(pattern, def...))
}

func (j *Json) GetInts(pattern string, def ...interface{}) []int {
	return conv.Ints(j.Get(pattern, def...))
}

func (j *Json) GetStrings(pattern string, def ...interface{}) []string {
	return conv.Strings(j.Get(pattern, def...))
}

func (j *Json) GetInterfaces(pattern string, def ...interface{}) []interface{} {
	return conv.Interfaces(j.Get(pattern, def...))
}

func (j *Json) GetTime(pattern string, format ...string) time.Time {
	return conv.Time(j.Get(pattern), format...)
}

func (j *Json) GetDuration(pattern string, def ...interface{}) time.Duration {
	return conv.Duration(j.Get(pattern, def...))
}

func (j *Json) GetVTime(pattern string, format ...string) *vtime.Time {
	return conv.VTime(j.Get(pattern), format...)
}

func (j *Json) Set(pattern string, value interface{}) error {
	return j.setValue(pattern, value, false)
}

func (j *Json) Remove(pattern string) error {
	return j.setValue(pattern, nil, true)
}

func (j *Json) Contains(pattern string) bool {
	return j.Get(pattern) != nil
}

func (j *Json) Len(pattern string) int {
	p := j.getPointerByPattern(pattern)
	if p != nil {
		switch (*p).(type) {
		case map[string]interface{}:
			return len((*p).(map[string]interface{}))
		case []interface{}:
			return len((*p).([]interface{}))
		default:
			return -1
		}
	}
	return -1
}

func (j *Json) Append(pattern string, value interface{}) error {
	p := j.getPointerByPattern(pattern)
	if p == nil {
		return j.Set(fmt.Sprintf("%s.0", pattern), value)
	}
	switch (*p).(type) {
	case []interface{}:
		return j.Set(fmt.Sprintf("%s.%d", pattern, len((*p).([]interface{}))), value)
	}
	return fmt.Errorf("invalid variable type of %s", pattern)
}

func (j *Json) GetStruct(pattern string, pointer interface{}, mapping ...map[string]string) error {
	return conv.Struct(j.Get(pattern), pointer, mapping...)
}

func (j *Json) GetStructDeep(pattern string, pointer interface{}, mapping ...map[string]string) error {
	return conv.StructDeep(j.Get(pattern), pointer, mapping...)
}

func (j *Json) GetStructs(pattern string, pointer interface{}, mapping ...map[string]string) error {
	return conv.Structs(j.Get(pattern), pointer, mapping...)
}

func (j *Json) GetStructsDeep(pattern string, pointer interface{}, mapping ...map[string]string) error {
	return conv.StructsDeep(j.Get(pattern), pointer, mapping...)
}

func (j *Json) GetScan(pattern string, pointer interface{}, mapping ...map[string]string) error {
	return conv.Scan(j.Get(pattern), pointer, mapping...)
}

func (j *Json) GetScanDeep(pattern string, pointer interface{}, mapping ...map[string]string) error {
	return conv.ScanDeep(j.Get(pattern), pointer, mapping...)
}

func (j *Json) GetMapToMap(pattern string, pointer interface{}, mapping ...map[string]string) error {
	return conv.MapToMap(j.Get(pattern), pointer, mapping...)
}

func (j *Json) GetMapToMapDeep(pattern string, pointer interface{}, mapping ...map[string]string) error {
	return conv.MapToMapDeep(j.Get(pattern), pointer, mapping...)
}

func (j *Json) GetMapToMaps(pattern string, pointer interface{}, mapping ...map[string]string) error {
	return conv.MapToMaps(j.Get(pattern), pointer, mapping...)
}

func (j *Json) GetMapToMapsDeep(pattern string, pointer interface{}, mapping ...map[string]string) error {
	return conv.MapToMapsDeep(j.Get(pattern), pointer, mapping...)
}

func (j *Json) ToMap() map[string]interface{} {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return conv.Map(*(j.p))
}

func (j *Json) ToArray() []interface{} {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return conv.Interfaces(*(j.p))
}

func (j *Json) ToStruct(pointer interface{}, mapping ...map[string]string) error {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return conv.Struct(*(j.p), pointer, mapping...)
}

func (j *Json) ToStructDeep(pointer interface{}, mapping ...map[string]string) error {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return conv.StructDeep(*(j.p), pointer, mapping...)
}

func (j *Json) ToStructs(pointer interface{}, mapping ...map[string]string) error {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return conv.Structs(*(j.p), pointer, mapping...)
}

func (j *Json) ToStructsDeep(pointer interface{}, mapping ...map[string]string) error {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return conv.StructsDeep(*(j.p), pointer, mapping...)
}

func (j *Json) ToScan(pointer interface{}, mapping ...map[string]string) error {
	return conv.Scan(*(j.p), pointer, mapping...)
}

func (j *Json) ToScanDeep(pointer interface{}, mapping ...map[string]string) error {
	return conv.ScanDeep(*(j.p), pointer, mapping...)
}

func (j *Json) ToMapToMap(pointer interface{}, mapping ...map[string]string) error {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return conv.MapToMap(*(j.p), pointer, mapping...)
}

func (j *Json) ToMapToMapDeep(pointer interface{}, mapping ...map[string]string) error {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return conv.MapToMapDeep(*(j.p), pointer, mapping...)
}

func (j *Json) ToMapToMaps(pointer interface{}, mapping ...map[string]string) error {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return conv.MapToMaps(*(j.p), pointer, mapping...)
}

func (j *Json) ToMapToMapsDeep(pointer interface{}, mapping ...map[string]string) error {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return conv.MapToMapsDeep(*(j.p), pointer, mapping...)
}

func (j *Json) Dump() {
	j.mu.RLock()
	defer j.mu.RUnlock()
	util.Dump(*j.p)
}

func (j *Json) Export() string {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return util.Export(*j.p)
}
