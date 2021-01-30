package conf

import (
	"encoding/json"
	"io/ioutil"
	file "utils/util"

	"gopkg.in/yaml.v2"
)

// 这里只留基础的获取文件数据的方法，提供json和yaml。 GetByKey和Get应该移到项目中。
// 获取配置文件可能不太需要包装的error方法，这类基础方法可能直接抛异常比较好。

func GetConfYaml(fileName string) (map[interface{}]interface{}, error) {
	c := make(map[interface{}]interface{})
	yamlFile, err := ioutil.ReadFile("conf/" + fileName + ".yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
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

func GetYamlByKey(key string) (interface{}, error) {
	m, err := GetConfYaml("conf")
	return m[key], err
}

func GetYaml(fileName, key string) (interface{}, error) {
	m, err := GetConfYaml(fileName)
	return m[key], err
}

//读取用户的配置文件
func Reload(ConfFilePath string, mod interface{}) {
	if confFileExists := file.PathExists(ConfFilePath); confFileExists != true {
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
