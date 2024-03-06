package console

import (
	"github.com/luyingjie/utils/encoding/base64"
	"github.com/luyingjie/utils/text/str"
)

func LoginDecode(encrypted, eid string) string {
	encrypteds := str.Split(encrypted, "@")
	if len(encrypteds) != 2 {
		return encrypted
	}
	signs, err1 := base64.DecodeString(encrypteds[0])
	if err1 != nil {
		return encrypted
	}
	pure := []rune(encrypteds[1])
	_eid := []rune(eid)
	signsArr := str.Split(string(signs), "")
	saltlen := len(_eid)
	b64 := ""
	b64u := make([]rune, len(pure))

	for i, v := range pure {
		var todel rune = 0
		if i < saltlen {
			todel = _eid[i]
		} else {
			todel = b64u[i-saltlen]
		}
		code := v*2 - todel
		if signsArr[i] == "1" {
			code += 1
		}
		if code != 0 {
			b64u[i] = code
			b64 += string(code)
		}
	}
	_decrypted, err := base64.DecodeToString(b64)
	if err != nil {
		return b64
	}

	return _decrypted
}

func LoginEncode(passwd, eid string) string {
	passwd = base64.EncodeString(passwd)
	_eid := []rune(eid)
	_passwd := []rune(passwd)
	passwdCount := len(passwd)
	eidCount := len(_eid)
	count := passwdCount
	if count < eidCount {
		count = eidCount
	}
	pure := make([]string, count)

	signs := ""

	for i := 0; i < count; i++ {
		var code rune
		var pcode rune
		var ecode rune
		if i < passwdCount {
			pcode = _passwd[i]
		} else {
			pcode = 0
		}
		if i < eidCount {
			ecode = _eid[i]
		} else {
			ecode = _passwd[i-eidCount]
		}
		if (pcode+ecode)%2 == 0 {
			pure[i] = "0"
		} else {
			pure[i] = "1"
			pcode -= 1
		}
		code = (pcode + ecode) / 2
		signs += string(code)
	}
	return base64.EncodeString(str.Join(pure, "")) + "@" + signs
}
