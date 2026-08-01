package utils

func InList[T comparable](key T, list []T) bool {
	for _, v := range list {
		if v == key {
			return true
		}
	}
	return false
}
