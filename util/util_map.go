package util

func MapCopy(data map[string]interface{}) (copy map[string]interface{}) {
	copy = make(map[string]interface{}, len(data))
	for k, v := range data {
		copy[k] = v
	}
	return
}

func MapContains(data map[string]interface{}, key string) (ok bool) {
	_, ok = data[key]
	return
}

func MapDelete(data map[string]interface{}, keys ...string) {
	if data == nil {
		return
	}
	for _, key := range keys {
		delete(data, key)
	}
}

func MapMerge(dst map[string]interface{}, src ...map[string]interface{}) {
	if dst == nil {
		return
	}
	for _, m := range src {
		for k, v := range m {
			dst[k] = v
		}
	}
}

func MapMergeCopy(src ...map[string]interface{}) (copy map[string]interface{}) {
	copy = make(map[string]interface{})
	for _, m := range src {
		for k, v := range m {
			copy[k] = v
		}
	}
	return
}

func MapPossibleItemByKey(data map[string]interface{}, key string) (foundKey string, foundValue interface{}) {
	if v, ok := data[key]; ok {
		return key, v
	}
	// Loop checking.
	for k, v := range data {
		if EqualFoldWithoutChars(k, key) {
			return k, v
		}
	}
	return "", nil
}

func MapContainsPossibleKey(data map[string]interface{}, key string) bool {
	if k, _ := MapPossibleItemByKey(data, key); k != "" {
		return true
	}
	return false
}

func MapOmitEmpty(data map[string]interface{}) {
	for k, v := range data {
		if IsEmpty(v) {
			delete(data, k)
		}
	}
}
