package conf

import (
	"encoding/json"
	"io/ioutil"
	myerr "utils/error"
	"utils/net/http"

	"gopkg.in/yaml.v2"
)

// 这里只留基础的获取文件数据的方法，提供json和yaml。 GetByKey和Get应该移到项目中。
// 获取配置文件可能不太需要包装的error方法，这类基础方法可能直接抛异常比较好。

func GetConfYaml(fileName string) map[interface{}]interface{} {
	c := make(map[interface{}]interface{})
	yamlFile, err := ioutil.ReadFile("conf/" + fileName + ".yaml")
	if err != nil {
		myerr.Try(3000, 3, "utils/conf/conf/GetConfYaml/ReadFile", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		myerr.Try(3000, 3, "utils/conf/conf/GetConfYaml/Unmarshal", err)
	}
	return c
}

func GetConf(fileName string) map[string]interface{} {
	c := make(map[string]interface{})
	Reload("conf/"+fileName+".json", &c)
	return c
}

func GetByKey(key string) interface{} {
	return GetConf("conf")[key]
}

func GetByKeyString(key string) map[string]string {
	o := GetByKey(key).(map[string]interface{})
	w := map[string]string{}
	for k, v := range o {
		w[k] = v.(string)
	}
	return w
}

func Get(fileName, key string) interface{} {
	return GetConf(fileName)[key]
}

func GetYamlByKey(key string) interface{} {
	return GetConfYaml("conf")[key]
}

func GetYaml(fileName, key string) interface{} {
	return GetConfYaml(fileName)[key]
}

//读取用户的配置文件
func Reload(ConfFilePath string, mod interface{}) {
	if confFileExists := http.PathExists(ConfFilePath); confFileExists != true {
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
