package moex

import (
	"bonds_calculator/internal/model/datastuct"
	"bonds_calculator/internal/utils"
	"bytes"
	"fmt"
	"io"
	"time"
)

// Structs for this response:
// https://iss.moex.com/iss/securities/$bondid/bondization.csv?limit=unlimited&iss.meta=off&iss.only=amortizations,coupons&amortizations.columns=$1&coupons.columns=$2
// $1 = amortdate,value
// $2 = coupondate,value

type (
	Bondization struct {
		Id            string
		Amortizations []Amortization
		Coupons       []Coupon
	}

	Amortization struct {
		Date  time.Time
		Value float64
	}

	Coupon struct {
		Date  time.Time
		Value datastuct.Optional[float64]
	}
)

func ParseBondization(id string, buf []byte) (Bondization, error) {
	reader := NewReader(bytes.NewReader(buf))

	header := ""

	amortizations := make([]Amortization, 0)
	coupons := make([]Coupon, 0)

	for {
		line, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return Bondization{}, err
		}

		if len(line) == 1 {
			if line[0] != "amortizations" && line[0] != "coupons" {
				return Bondization{}, fmt.Errorf("invalid header %s", line[0])
			}

			header = line[0]

			continue
		}

		item, err := tryParseItem(line)

		if err != nil {
			return Bondization{}, err
		}

		if header == "amortizations" {
			value, exist := item.Value.Get()

			if !exist {
				return Bondization{}, fmt.Errorf("amortizations value must be not null")
			}

			amortizations = append(amortizations, Amortization{
				Date:  item.Date,
				Value: value,
			})
		} else {
			coupons = append(coupons, Coupon(item))
		}
	}

	return Bondization{
		Id:            id,
		Amortizations: amortizations,
		Coupons:       coupons,
	}, nil
}

type commonItem struct {
	Date  time.Time
	Value datastuct.Optional[float64]
}

func tryParseItem(line []string) (commonItem, error) {
	if len(line) != 2 {
		return commonItem{}, fmt.Errorf("wrong Amortization data line len %d", len(line))
	}

	date, err := time.Parse("2006-01-02", line[0])
	value, err := utils.ParseOptionalFloat64(line[1])

	if err != nil {
		return commonItem{}, fmt.Errorf("cannot parse Bondization item %v", err)
	}

	return commonItem{
		Date:  date,
		Value: value,
	}, nil
}
