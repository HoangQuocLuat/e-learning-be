package utilities

func SetString(value string) *string {
	return &value
}

func SetBool(value bool) *bool {
	return &value
}

func StringIntArray(str string, arr []string) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] == str {
			return true
		}
	}

	return false
}
