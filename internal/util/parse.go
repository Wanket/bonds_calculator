package util

import (
	"fmt"
	"time"
)

func ParseMoexDate(str string) (time.Time, error) {
	if str == "" {
		return time.Time{}, nil
	}

	res, err := time.Parse("2006-01-02", str)
	if err != nil {
		if str != "0000-00-00" {
			return time.Time{}, fmt.Errorf("cannot parse moex date %w", err)
		}

		return time.Time{}, nil
	}

	return res, nil
}
