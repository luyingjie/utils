package util

func SliceCopy(data []interface{}) []interface{} {
	newData := make([]interface{}, len(data))
	copy(newData, data)
	return newData
}

func SliceDelete(data []interface{}, index int) (newSlice []interface{}) {
	if index < 0 || index >= len(data) {
		return data
	}
	if index == 0 {
		return data[1:]
	} else if index == len(data)-1 {
		return data[:index]
	}

	return append(data[:index], data[index+1:]...)
}
