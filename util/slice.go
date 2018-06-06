package util

import "strings"

func InSlice(str string, values []string) bool {
	for _, v := range values {
		if str == v {
			return true
		}
	}
	return false
}

func InSliceContains(fullString string, candidates []string) bool {
	for _, v := range candidates {
		if strings.Contains(fullString, v) {
			return true
		}
	}
	return false
}
