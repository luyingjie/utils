package iaas

import (
	"fmt"
	"testing"
)

var conf map[string]interface{} = map[string]interface{}{
	"host":                "api.pekdemo.com",
	"port":                "7777",
	"protocol":            "http",
	"console_uri":         "/iaas/",
	"console_key_id":      "OXIQRFHUWMQWHPYOPATL",
	"console_secrect_key": "0BH3KWpluMikGedYAo0sJi0vLmbOGEfeZ6Jq1eZt",
	"zone":                "pekdemo1",
}

var params map[string]interface{} = map[string]interface{}{
	"action": "DescribeInstances",
	"instances": []string{
		"i-e8mwsqlo",
		"i-y0mp12lh",
		"i-klyrrs9h",
		"i-gz2evb0o",
		"i-gbdxw2p0",
		"i-p025x7u6",
		"i-wzwu5vhe",
		"i-p20u1bqd",
		"i-wox7vuvn",
		"i-ui2kf31j",
		// "i-7kyomx5p",
		// "i-bhp60yhc",
		// "i-l0136k6y",
		// "i-3s3f63yt",
		// "i-836rqqs6",
		// "i-cnmeuvyx",
		// "i-rtj681wb",
		// "i-bglqc980",
		// "i-qryv28be",
		// "i-1b3938lq",
	},
	// "instances.1":  "i-e8mwsqlo",
	// "instances.2":  "i-y0mp12lh",
	// "instances.3":  "i-klyrrs9h",
	// "instances.4":  "i-gz2evb0o",
	// "instances.5":  "i-gbdxw2p0",
	// "instances.6":  "i-p025x7u6",
	// "instances.7":  "i-wzwu5vhe",
	// "instances.8":  "i-p20u1bqd",
	// "instances.9":  "i-wox7vuvn",
	// "instances.10": "i-ui2kf31j",
	// "instances.11": "i-7kyomx5p",
	// "instances.12": "i-bhp60yhc",
	"verbose": 0,
	"zone":    "pekdemo1",
}

func TestEncoding(t *testing.T) {

	ot, _ := Send("GET", params, conf)
	fmt.Println(ot)
}
