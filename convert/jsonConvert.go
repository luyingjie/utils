package convert

import (
	"encoding/json"
	"utils/error"
	// "github.com/bitly/go-simplejson"
)

// func InterfaceToMap(i interface{}) (map[string]string){
// 	m := i.(map[string]interface{})

// 	for k, v := range m {
// 		switch vv := v.(type) {
// 		case string:
// 			fmt.Println(k, "is string", vv)
// 		case int:
// 			fmt.Println(k, "is int", vv)
// 		case float64:
// 			fmt.Println(k,"is float64",vv)
// 		case []interface{}:
// 			fmt.Println(k, "is an array:")
// 			for i, u := range vv {
// 				fmt.Println(i, u)
// 			}
// 		default:
// 			fmt.Println(k, "is of a type I don't know how to handle")
// 		}
// 	}

// 	return nil
// }

//将一个map[string]string序列化成json的字符串
func MapToString(data map[string]string) string {
	str, err := json.Marshal(data)
	if err != nil {
		error.Try(5000, 3, err)
	}
	return string(str)
}

//将一个字串转成map
func ByteToMap(data []byte) map[string]string {
	var model interface{}
	json.Unmarshal(data, &model)
	requestModel := make(map[string]string)
	for key, value := range model.(map[string]interface{}) {
		requestModel[key] = value.(string)
	}
	return requestModel
}

func ByteToMapInterface(data []byte) map[string]interface{} {
	var model interface{}
	json.Unmarshal(data, &model)
	// requestModel := make(map[string]string)
	// for key, value := range model.(map[string]interface{}) {
	// 	requestModel[key] = value
	// }
	// return requestModel
	return model.(map[string]interface{})
}
