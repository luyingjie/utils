package console

import (
	"fmt"
	"testing"
)

func TestEncoding(t *testing.T) {
	// eid   to base64   32 -> 44
	// lmiXQdVlvfzfddhPZkjSTXarcvmjnkqc
	// bG1pWFFkVmx2ZnpmZGRoUFpralNUWGFyY3Ztam5rcWM=

	// encrypted
	// admin@staging.com Zhu88jie
	// MTAxMTAwMDExMTExMDAwMDAxMDEwMDEwMTAxMDAxMTE=@amhDPT_nk^gQ224(-55)*,091;657581

	// MTAxMTAwMDExMTExMDAwMDAxMDEwMDEwMTAxMDAxMTE=  to base64
	// 10110001111100000101001010100111

	// ASCII

	// var sss rune = 'a'
	// var ssss int = 87
	// fmt.Println(int(sss))
	// fmt.Println(string(rune(ssss)))

	// name := "admin@staging.com"
	// passwd := "Zhu88jie"

	// 长度大道32位后b64后会成为44位
	eid := "lmiXQdVlvfzfddhPZkjSTXarcvmjnkqc"
	passwd := "Zhdyrutjdkfiririririririririri"

	decrypted := LoginEncode(passwd, eid)

	fmt.Println(decrypted)
	fmt.Println(LoginDecode(decrypted, eid))
}
