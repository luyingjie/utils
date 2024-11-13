package aes

import (
	"fmt"
	"testing"
)

func TestAES(t *testing.T) {
	orig := "Zhu88jie"
	key1 := "ZG9tYWluX2xvZ2lu" //"0123456789012345"
	fmt.Println("原文：", orig)
	encryptCode := AesEncrypt(orig, key1)
	fmt.Println("密文：", encryptCode)
	// decryptCode := AesDecrypt(encryptCode, key1)
	decryptCode := AesDecrypt("GZA7mxUw8q+zCso8vjNqk0lU93/unlzxyG0rDAtu+4Fp8c8gBrASkaZJpecN/KbjUiLQhhMSZExS/iEpdxAs7qPv/BmGYcWq0il6wFm9OtM=", key1)
	fmt.Println("解密结果：", decryptCode)
}
