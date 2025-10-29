package type_converter

// MapToArray converts a map of pointers to a slice of values.
// If the input map is nil, it returns nil.
func MapToArray[T any](m map[string]*T) []T {
	if m == nil {
		return nil
	}

	arr := make([]T, 0, len(m))
	for _, val := range m {
		arr = append(arr, *val)
	}

	return arr
}
