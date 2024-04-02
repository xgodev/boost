package strings

// SliceContains receives a slice of string and checks if contains str
func SliceContains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}

	return false
}

// SliceNotContains receives a slice of string and checks if contains str
func SliceNotContains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return false
		}
	}

	return true
}

// SliceContainsOneOf receives a slice of string and checks if contains one of string
func SliceContainsOneOf(slice []string, strs []string) bool {
	for _, s := range slice {
		for _, c := range strs {
			if s == c {
				return true
			}
		}
	}

	return false
}

// SliceNotContainsOneOf receives a slice of string and checks if contains one of string
func SliceNotContainsOneOf(slice []string, strs []string) bool {
	for _, s := range slice {
		for _, c := range strs {
			if s == c {
				return false
			}
		}
	}

	return true
}

// SliceContainsAll receives a slice of string and checks if contains all of string slice
func SliceContainsAll(slice []string, strs []string) bool {
	for _, s := range slice {
		for _, c := range strs {
			if s != c {
				return false
			}
		}
	}

	return true
}
