package conf

import (
	"encoding/json"
	"io/ioutil"
	"utils/base"
	"utils/error"

	"gopkg.in/yaml.v2"
)

// 这里只留基础的获取文件数据的方法，提供json和yaml。 GetByKey和Get应该移到项目中。
// 获取配置文件可能不太需要包装的error方法，这类基础方法可能直接抛异常比较好。

func GetConfYaml(fileName string) map[interface{}]interface{} {
	c := make(map[interface{}]interface{})
	yamlFile, err := ioutil.ReadFile("conf/" + fileName + ".yaml")
	if err != nil {
		error.TryError(err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		error.TryError(err)
	}
	return c
}

func GetConf(fileName string) map[string]interface{} {
	c := make(map[string]interface{})
	Reload("conf/" + fileName + ".json", &c)
	return c
}

func GetByKey(key string) interface{} {
	return GetConf("conf")[key]
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
	if confFileExists, _ := base.PathExists(ConfFilePath); confFileExists != true {
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
