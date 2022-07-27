package util

func SliceToMapBy[K comparable, V any, F func(V) K](values []V, f F) map[K]V {
	result := make(map[K]V, len(values))

	for _, value := range values {
		result[f(value)] = value
	}

	return result
}

func AnyOf[T any](arr []T, fun func(T) bool) bool {
	for _, item := range arr {
		if fun(item) {
			return true
		}
	}

	return false
}
