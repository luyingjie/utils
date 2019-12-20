package crypto

import (
	"fmt"
	"testing"
)

func TestAES(t *testing.T) {
	orig := "hello world"
	key1 := "luke@yunify.com0" //"0123456789012345"
	fmt.Println("原文：", orig)
	encryptCode := AesEncrypt(orig, key1)
	fmt.Println("密文：", encryptCode)
	decryptCode := AesDecrypt(encryptCode, key1)
	fmt.Println("解密结果：", decryptCode)
}
