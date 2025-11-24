package datatype

// IntToBool 将整数转换为布尔值
func IntToBool(i int) bool {
	return i != 0
}

// ContainsInt 检查整数是否在切片中
func ContainsInt(slice []int, num int) bool {
	for _, v := range slice {
		if v == num {
			return true
		}
	}
	return false
}
