package rand

import (
	"fmt"
	"testing"
)

func TestRand(t *testing.T) {
	v := Str("1234567890", 6)
	fmt.Println(v)
	fmt.Println(Str("1234567890", 6))
	fmt.Println(Str("1234567890", 6))
	fmt.Println(Str("1234567890", 6))
	fmt.Println(Str("1234567890", 6))
}
