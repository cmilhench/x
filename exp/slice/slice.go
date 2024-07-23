package slice

// Take returns the first n elements of the given slice. If there are not
// enough elements in the slice, the whole slice is returned.
func Take[A any](slice []A, n int) []A {
	if n > len(slice) {
		return slice
	}
	return slice[:n]
}

// Unique returns a new slice containing only the unique elements of the given
// slice. The order of the elements is preserved.
func Unique[A comparable](slice []A) []A {
	seen := make(map[A]bool, len(slice))
	unique := []A{}
	for _, s := range slice {
		if !seen[s] {
			seen[s] = true
			unique = append(unique, s)
		}
	}
	return unique
}
