package str

import (
	"strings"
	"utils/os/conv"
)

func CompareVersion(a, b string) int {
	if a[0] == 'v' {
		a = a[1:]
	}
	if b[0] == 'v' {
		b = b[1:]
	}
	var (
		array1 = strings.Split(a, ".")
		array2 = strings.Split(b, ".")
		diff   = 0
	)
	diff = len(array2) - len(array1)
	for i := 0; i < diff; i++ {
		array1 = append(array1, "0")
	}
	diff = len(array1) - len(array2)
	for i := 0; i < diff; i++ {
		array2 = append(array2, "0")
	}
	v1 := 0
	v2 := 0
	for i := 0; i < len(array1); i++ {
		v1 = conv.Int(array1[i])
		v2 = conv.Int(array2[i])
		if v1 > v2 {
			return 1
		}
		if v1 < v2 {
			return -1
		}
	}
	return 0
}

func CompareVersionGo(a, b string) int {
	if a[0] == 'v' {
		a = a[1:]
	}
	if b[0] == 'v' {
		b = b[1:]
	}
	if Count(a, "-") > 1 {
		if i := PosR(a, "-"); i > 0 {
			a = a[:i]
		}
	}
	if Count(b, "-") > 1 {
		if i := PosR(b, "-"); i > 0 {
			b = b[:i]
		}
	}
	if i := Pos(a, "+"); i > 0 {
		a = a[:i]
	}
	if i := Pos(b, "+"); i > 0 {
		b = b[:i]
	}
	a = Replace(a, "-", ".")
	b = Replace(b, "-", ".")
	var (
		array1 = strings.Split(a, ".")
		array2 = strings.Split(b, ".")
		diff   = 0
	)
	// Specially in Golang:
	// "v1.12.2-0.20200413154443-b17e3a6804fa" < "v1.12.2"
	if len(array1) > 3 && len(array2) <= 3 {
		return -1
	}
	if len(array1) <= 3 && len(array2) > 3 {
		return 1
	}
	diff = len(array2) - len(array1)
	for i := 0; i < diff; i++ {
		array1 = append(array1, "0")
	}
	diff = len(array1) - len(array2)
	for i := 0; i < diff; i++ {
		array2 = append(array2, "0")
	}
	v1 := 0
	v2 := 0
	for i := 0; i < len(array1); i++ {
		v1 = conv.Int(array1[i])
		v2 = conv.Int(array2[i])
		if v1 > v2 {
			return 1
		}
		if v1 < v2 {
			return -1
		}
	}
	return 0
}
