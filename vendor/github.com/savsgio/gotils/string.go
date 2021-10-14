package gotils

// StringSliceIndexOf returns index position in slice from given string
// If value is -1, the string does not found.
func StringSliceIndexOf(vs []string, s string) int {
	for i, v := range vs {
		if v == s {
			return i
		}
	}

	return -1
}

// StringSliceInclude returns true or false if given string is in slice.
func StringSliceInclude(vs []string, t string) bool {
	return StringSliceIndexOf(vs, t) >= 0
}
