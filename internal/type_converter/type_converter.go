package type_converter

// MapValueToArray converts a map of pointers to a slice of values.
// If the input map is nil, it returns nil.
func MapValueToArray[T any](m map[string]*T) []T {
	if m == nil {
		return nil
	}

	arr := make([]T, 0, len(m))
	for _, val := range m {
		arr = append(arr, *val)
	}

	return arr
}

// MapKeyToArray converts a map to a slice of its keys.
// If the input map is nil, it returns nil.
func MapKeyToArray[T any](m map[string]T) []string {
	if m == nil {
		return nil
	}

	arr := make([]string, 0, len(m))
	for k := range m {
		arr = append(arr, k)
	}

	return arr
}
