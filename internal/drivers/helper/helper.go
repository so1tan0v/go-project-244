package helper

func MergeMaps[K comparable, V any](maps ...map[K]V) map[K]V {
	merged := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			merged[k] = v
		}
	}
	return merged
}

func MergeSlices[T any](slices ...[]T) []T {
	var merged []T
	for _, s := range slices {
		merged = append(merged, s...)
	}
	return merged
}

func FilterSlices[TKey any, T []TKey](slice T, fn func(k int, i TKey) bool) T {
	result := make(T, 0)

	for k, s := range slice {
		if fn(k, s) {
			result = append(result, s)
		}
	}

	return result
}
