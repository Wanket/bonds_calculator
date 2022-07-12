package util

import (
	"bonds_calculator/internal/model/datastuct"
	"strconv"
)

func ParseOptionalFloat64(str string) (datastuct.Optional[float64], error) {
	if str == "" {
		return datastuct.Optional[float64]{}, nil
	}

	res, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return datastuct.Optional[float64]{}, err
	}

	return datastuct.NewOptional(res), nil
}
