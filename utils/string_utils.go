package utils

// IsStringSliceContains -- check slice contain string
func IsStringSliceContains(stringSlice []string, searchString string) bool {
	for _, value := range stringSlice {
		if value == searchString {
			return true
		}
	}
	return false
}