package conv

type apiString interface {
	String() string
}

type apiError interface {
	Error() string
}

type apiInterfaces interface {
	Interfaces() []interface{}
}

type apiFloats interface {
	Floats() []float64
}

type apiInts interface {
	Ints() []int
}

type apiStrings interface {
	Strings() []string
}

type apiUints interface {
	Uints() []uint
}

type apiMapStrAny interface {
	MapStrAny() map[string]interface{}
}

type apiUnmarshalValue interface {
	UnmarshalValue(interface{}) error
}

type apiSet interface {
	Set(value interface{}) (old interface{})
}
