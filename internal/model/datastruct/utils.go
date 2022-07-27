package datastruct

import (
	"fmt"
	"strconv"
)

func ParseOptionalFloat64(str string) (Optional[float64], error) {
	if str == "" {
		return Optional[float64]{}, nil
	}

	res, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return Optional[float64]{}, fmt.Errorf("cannot parse float: %w", err)
	}

	return NewOptional(res), nil
}
