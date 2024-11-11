package rand

import (
	"fmt"
	"testing"
)

func TestRand(t *testing.T) {
	v := Intn(999999)
	fmt.Println(v)
	fmt.Println(Intn(999999))
	fmt.Println(Intn(999999))
	fmt.Println(Intn(999999))
	fmt.Println(Intn(999999))

	v1 := Str("1234567890", 6)
	fmt.Println(v1)
	fmt.Println(Str("1234567890", 6))
	fmt.Println(Str("1234567890", 6))
	fmt.Println(Str("1234567890", 6))
	fmt.Println(Str("1234567890", 6))
}
