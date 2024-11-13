package xml

import (
	"strings"

	"github.com/luyingjie/utils/text/regex"

	"github.com/luyingjie/utils/encoding/charset"

	"github.com/clbanning/mxj"
)

func Decode(content []byte) (map[string]interface{}, error) {
	res, err := convert(content)
	if err != nil {
		return nil, err
	}
	return mxj.NewMapXml(res)
}

func DecodeWithoutRoot(content []byte) (map[string]interface{}, error) {
	res, err := convert(content)
	if err != nil {
		return nil, err
	}
	m, err := mxj.NewMapXml(res)
	if err != nil {
		return nil, err
	}
	for _, v := range m {
		if r, ok := v.(map[string]interface{}); ok {
			return r, nil
		}
	}
	return m, nil
}

func Encode(m map[string]interface{}, rootTag ...string) ([]byte, error) {
	return mxj.Map(m).Xml(rootTag...)
}

func EncodeWithIndent(m map[string]interface{}, rootTag ...string) ([]byte, error) {
	return mxj.Map(m).XmlIndent("", "\t", rootTag...)
}

func ToJson(content []byte) ([]byte, error) {
	res, err := convert(content)
	if err != nil {
		return nil, err
	}
	mv, err := mxj.NewMapXml(res)
	if err == nil {
		return mv.Json()
	} else {
		return nil, err
	}
}

func convert(xml []byte) (res []byte, err error) {
	patten := `<\?xml.*encoding\s*=\s*['|"](.*?)['|"].*\?>`
	matchStr, err := regex.MatchString(patten, string(xml))
	if err != nil {
		return nil, err
	}
	xmlEncode := "UTF-8"
	if len(matchStr) == 2 {
		xmlEncode = matchStr[1]
	}
	xmlEncode = strings.ToUpper(xmlEncode)
	res, err = regex.Replace(patten, []byte(""), xml)
	if err != nil {
		return nil, err
	}
	if xmlEncode != "UTF-8" && xmlEncode != "UTF8" {
		dst, err := charset.Convert("UTF-8", xmlEncode, string(res))
		if err != nil {
			return nil, err
		}
		res = []byte(dst)
	}
	return res, nil
}
