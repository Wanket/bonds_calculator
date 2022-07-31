package moex

import (
	"bonds_calculator/internal/model/datastruct"
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/exp/slices"
	"io"
	"math"
	"time"
)

type (
	Bondization struct {
		ID            string
		Amortizations []Amortization
		Coupons       []Coupon
	}

	Amortization struct {
		Date  time.Time
		Value float64
	}

	Coupon struct {
		Date  time.Time
		Value datastruct.Optional[float64]
	}
)

var (
	errEmptyBondizationID   = errors.New("empty bondization id")
	errLastAmortizationDate = errors.New("last amortization date must be equal to end date")

	errEmptyAmortizations    = errors.New("amortizations is empty")
	errAmortizationsUnsorted = errors.New("amortizations must be sorted")
	errAmortizationsValue    = errors.New("amortizations value must be > 0")

	errEmptyCoupons    = errors.New("coupons is empty")
	errCouponsUnsorted = errors.New("coupons must be sorted")
	errCouponsValue    = errors.New("coupons value must be > 0")
	errCouponsNan      = errors.New("coupons value must be not NaN")

	errAmortizationValueIsNull = errors.New("amortization value is null")

	errWrongAmortizationDataLineLength = errors.New("wrong amortization data line length")
)

func (bondization *Bondization) IsValid(endDate time.Time) error {
	if bondization.ID == "" {
		return errEmptyBondizationID
	}

	if err := CheckAmortizations(bondization.Amortizations); err != nil {
		return err
	}

	if !endDate.IsZero() && bondization.Amortizations[len(bondization.Amortizations)-1].Date != endDate {
		return errLastAmortizationDate
	}

	if err := CheckCoupons(bondization.Coupons); err != nil {
		return err
	}

	return nil
}

func CheckAmortizations(amortizations []Amortization) error {
	if len(amortizations) == 0 {
		return errEmptyAmortizations
	}

	if !slices.IsSortedFunc(amortizations, func(left, right Amortization) bool {
		return left.Date.Before(right.Date)
	}) {
		return errAmortizationsUnsorted
	}

	for _, amortization := range amortizations {
		if amortization.Value < 0 {
			return errAmortizationsValue
		}
	}

	return nil
}

func CheckCoupons(coupons []Coupon) error {
	if len(coupons) == 0 {
		return errEmptyCoupons
	}

	if !slices.IsSortedFunc(coupons, func(left, right Coupon) bool {
		return left.Date.Before(right.Date)
	}) {
		return errCouponsUnsorted
	}

	for _, coupon := range coupons {
		if value, exist := coupon.Value.Get(); exist {
			if value <= 0 {
				return errCouponsValue
			}

			if math.IsNaN(value) {
				return errCouponsNan
			}
		}
	}

	return nil
}

func ParseBondization(bondID string, buf []byte) (Bondization, error) {
	reader := NewReader(bytes.NewReader(buf))

	header := ""

	amortizations := make([]Amortization, 0)
	coupons := make([]Coupon, 0)

	for line, err := reader.Read(); !errors.Is(err, io.EOF); line, err = reader.Read() {
		if err != nil {
			return Bondization{}, fmt.Errorf("cannot read line: %w", err)
		}

		if len(line) == 1 {
			if line[0] != "amortizations" && line[0] != "coupons" {
				return Bondization{}, fmt.Errorf("ParseBondization: %w: %s", errInvalidHeader, line[0])
			}

			header = line[0]

			continue
		}

		err := parseItemToArray(header, line, &amortizations, &coupons)

		if err != nil {
			return Bondization{}, err
		}
	}

	if len(amortizations) == 0 {
		amortizations = append(amortizations, Amortization{})
	}

	if len(coupons) == 0 {
		coupons = append(coupons, Coupon{})
	}

	return Bondization{
		ID:            bondID,
		Amortizations: amortizations,
		Coupons:       coupons,
	}, nil
}

type commonItem struct {
	Date  time.Time
	Value datastruct.Optional[float64]
}

func parseItemToArray(header string, line []string, amortizations *[]Amortization, coupons *[]Coupon) error {
	item, err := tryParseItem(line)

	if err != nil {
		return err
	}

	if header == "amortizations" {
		value, exist := item.Value.Get()
		if !exist {
			return errAmortizationValueIsNull
		}

		*amortizations = append(*amortizations, Amortization{
			Date:  item.Date,
			Value: value,
		})

		return nil
	}

	*coupons = append(*coupons, Coupon(item))

	return nil
}

func tryParseItem(line []string) (commonItem, error) {
	const commonItemLineSize = 2

	if len(line) != commonItemLineSize {
		return commonItem{}, fmt.Errorf("tryParseItem: %w %d", errWrongAmortizationDataLineLength, len(line))
	}

	date, err := time.Parse("2006-01-02", line[0])
	if err != nil {
		return commonItem{}, fmt.Errorf("cannot parse Bondization item %w", err)
	}

	value, err := datastruct.ParseOptionalFloat64(line[1])
	if err != nil {
		return commonItem{}, fmt.Errorf("cannot parse Bondization item %w", err)
	}

	return commonItem{
		Date:  date,
		Value: value,
	}, nil
}
