package utils

import (
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
