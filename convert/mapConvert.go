package convert

// InterfaceToString : 处理Yaml文件读出来的结果的不确定性，其他类型可以直接输出，需要处理最外层的map[interface{}]interface{}为map[string]interface{}
func InterfaceToString(p interface{}) interface{} {
	// map[interface{}]interface{}
	switch p.(type) {
	case map[interface{}]interface{}:
		s := map[string]interface{}{}
		for key, val := range p.(map[interface{}]interface{}) {
			s[key.(string)] = val
		}
		return s
	}
	return p
}
