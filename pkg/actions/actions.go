package actions

import "strings"

func Marshal(pending string) string {
	if pending == "" {
		return ""
	}

	temp := strings.Split(pending, "_")
	var s string
	for _, v := range temp {
		chv := []rune(v)
		if len(chv) > 0 {
			if chv[0] >= 'a' && chv[0] <= 'z' { //首字母大写
				chv[0] -= 32
			}
			s += string(chv)
		}
	}
	return s
}

func IsStringInSlice(item string, slice []string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func IsIntInSlice(item int, slice []int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func IntToBool(i int) bool {
	return i != 0
}
