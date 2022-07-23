package boss2

import (
	"fmt"
	"testing"
)

var conf map[string]interface{} = map[string]interface{}{
	"host":                "boss.testing.com", //"192.168.27.5"
	"port":                "80",
	"protocol":            "http",
	"console_uri":         "/boss2/",
	"console_key_id":      "AKEQCMDFBXBSXWUZIISX",
	"console_secrect_key": "jt3grf70HIynwq5IJtocOI1xp36hAe91maz72r4p",
	"zone":                "testing1a",
}

var params map[string]interface{} = map[string]interface{}{
	"action": "Boss2DescribeBots",
	"zone":   "testing1a",
	"limit":  1,
}

func TestCheck(t *testing.T) {
	// sig, s, _, _ := Signature("/bos2/", conf["console_uri"].(string), conf["console_key_id"].(string), params)
	sig, s, _, _ := SignatureStr("/boss2/", "AKEQCMDFBXBSXWUZIISX", "jt3grf70HIynwq5IJtocOI1xp36hAe91maz72r4p", "Boss2DescribeBots", "{\"action\": \"Boss2DescribeBots\", \"zone\": \"testing\", \"limit\": 1}")
	fmt.Println(sig)
	fmt.Println(s)
}

func TestEncoding(t *testing.T) {

	ot, _ := Send(params, conf)
	fmt.Println(ot)
}
