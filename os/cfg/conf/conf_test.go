package conf

 import (
     "testing"
     "fmt"
 )



 func TestMain(t *testing.T) {
	 fmt.Println("开始")
	 var json interface{}
	 Reload("test.json", &json)
	 
	 fmt.Println(json.(map[string]interface{})["key1"].(string))

	 var json2 map[string]interface{}
	 Reload("test.json", &json2)
	 
	 fmt.Println(json2["key2"].(float64))

     fmt.Println("结束")
 }
