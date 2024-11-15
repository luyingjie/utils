package util

import (
	"encoding/json"
	"io/ioutil"
)

func ReloadJson(ConfFilePath string, mod interface{}) {
	if confFileExists := PathExists(ConfFilePath); confFileExists != true {
		return
	}

	data, err := ioutil.ReadFile(ConfFilePath)
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	err = json.Unmarshal(data, mod)
	if err != nil {
		panic(err)
	}
}
