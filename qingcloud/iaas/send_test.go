package iaas

import (
	"encoding/json"
	"fmt"
	"testing"
)

var conf map[string]interface{} = map[string]interface{}{
	"host":                "192.168.10.131",
	"port":                "7777",
	"protocol":            "http",
	"console_uri":         "/iaas/",
	"console_key_id":      "LZZYPNKNWUZENPTMCDSK",
	"console_secrect_key": "yXbvLkEHnKZtbHwBYPdovk7vLf5vYJqjiZWvpAPO",
	"zone":                "devops1a",
}

func TestDescribeZones(t *testing.T) {
	var params map[string]interface{} = map[string]interface{}{
		"action": "DescribeZones",
	}
	ot, _ := Send("GET", params, conf)
	os, _ := json.Marshal(ot)
	fmt.Println(string(os))
}

func TestDescribeAccessKeys(t *testing.T) {
	var params map[string]interface{} = map[string]interface{}{
		"action":        "DescribeAccessKeys",
		"access_keys.1": "LZZYPNKNWUZENPTMCDSK",
		"zone":          "devops1a",
	}
	ot, _ := Send("GET", params, conf)
	os, _ := json.Marshal(ot)
	fmt.Println(string(os))
}
