package vmap

type (
	Map     = AnyAnyMap
	HashMap = AnyAnyMap
)

func New(safe ...bool) *Map {
	return NewAnyAnyMap(safe...)
}

func NewFrom(data map[interface{}]interface{}, safe ...bool) *Map {
	return NewAnyAnyMapFrom(data, safe...)
}

func NewHashMap(safe ...bool) *Map {
	return NewAnyAnyMap(safe...)
}

func NewHashMapFrom(data map[interface{}]interface{}, safe ...bool) *Map {
	return NewAnyAnyMapFrom(data, safe...)
}
