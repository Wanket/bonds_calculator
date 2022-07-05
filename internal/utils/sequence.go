package utils

import (
	"bonds_calculator/internal/model/datastuct"
	"golang.org/x/exp/constraints"
)

func SumBy[T any, R constraints.Ordered, F func(T) R](arr []T, fun F) R {
	var sum R
	for _, item := range arr {
		sum += fun(item)
	}

	return sum
}

func AvgBy[T any, R constraints.Integer | constraints.Float, F func(T) R](arr []T, fun F) R {
	return SumBy(arr, fun) / R(len(arr))
}

func MapToSlice[K comparable, V any](m map[K]V) []datastuct.Pair[K, V] {
	result := make([]datastuct.Pair[K, V], 0, len(m))

	for k, v := range m {
		result = append(result, datastuct.Pair[K, V]{k, v})
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
