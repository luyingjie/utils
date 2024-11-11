package xml

import (
	"fmt"
	"testing"
)

func Test_DecodeWitoutRoot(t *testing.T) {
	content := `
	<?xml version="1.0" encoding="UTF-8"?><doc><username>johngcn</username><password1>123456</password1><password2>123456</password2></doc>
	`
	// Decode
	m, _ := DecodeWithoutRoot([]byte(content))

	fmt.Println(m)
	fmt.Println(m["username"])
	fmt.Println(m["password1"])
	fmt.Println(m["password2"])

}
