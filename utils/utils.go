package utils

func ExistsInArray(arr []string, search string) bool {
	for _, str := range arr {
		if str == search {
			return true
		}
	}
	return false
}