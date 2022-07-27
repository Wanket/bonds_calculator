package moex

import (
	"bonds_calculator/internal/model/datastruct"
	"bonds_calculator/internal/util"
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"io"
	"strconv"
	"time"
)

type (
	Bond struct {
		ID string
		SecurityPart
		MarketDataPart
	}

	SecurityPart struct {
		ShortName        string
		Coupon           float64
		NextCoupon       time.Time
		AccCoupon        float64
		PrevPricePercent float64
		Value            float64
		CouponPeriod     uint
		PriceStep        float64
		CouponPercent    datastruct.Optional[float64]
		EndDate          time.Time
		Currency         string
	}

	MarketDataPart struct {
		CurrentPricePercent float64
	}
)

var (
	errBondIDIsEmpty                       = errors.New("bond id is empty")
	errBondShortNameIsEmpty                = errors.New("bond short name is empty")
	errBondCouponLessThanZero              = errors.New("bond coupon is < 0")
	errBondAccCouponLessThanZero           = errors.New("bond acc coupon is < 0")
	errBondValueLessThanZero               = errors.New("bond value is <= 0")
	errBondCouponPeriodXorNextCouponIsZero = errors.New("bond coupon period is 0 xor next coupon is zero")
	errBondPrevPricePercentLessThanZero    = errors.New("bond prev price percent is <= 0")
)

func (bond *Bond) IsValid() error {
	if bond.ID == "" {
		return errBondIDIsEmpty
	}

	if bond.ShortName == "" {
		return errBondShortNameIsEmpty
	}

	if bond.Coupon < 0 {
		return errBondCouponLessThanZero
	}

	if bond.AccCoupon < 0 {
		return errBondAccCouponLessThanZero
	}

	if bond.Value <= 0 {
		return errBondValueLessThanZero
	}

	zeroCouponPeriod, zeroNextCoupon := bond.CouponPeriod == 0, bond.NextCoupon.IsZero()
	if zeroCouponPeriod != zeroNextCoupon {
		return errBondCouponPeriodXorNextCouponIsZero
	}

	if bond.PriceStep <= 0 {
		return errBondPrevPricePercentLessThanZero
	}

	return nil
}

func (bond *Bond) AbsoluteCurrentPrice() float64 {
	const percentMultiplier = 100

	return bond.CurrentPricePercent * bond.Value / percentMultiplier
}

func ParseBondsCp1251(buf []byte) ([]Bond, error) {
	decoded, err := charmap.Windows1251.NewDecoder().Bytes(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot decode cp1251: %w", err)
	}

	return ParseBonds(decoded)
}

var errInvalidHeader = errors.New("invalid header")

func ParseBonds(buf []byte) ([]Bond, error) {
	reader := NewReader(bytes.NewReader(buf))

	header := ""

	resultMap := make(map[string]Bond)

	for line, err := reader.Read(); !errors.Is(err, io.EOF); line, err = reader.Read() {
		if err != nil {
			return nil, fmt.Errorf("cannot read line: %w", err)
		}

		if len(line) == 1 {
			if line[0] != "marketdata" && line[0] != "securities" {
				return nil, fmt.Errorf("ParseBonds: %w: %s", errInvalidHeader, line[0])
			}

			header = line[0]

			continue
		}

		if header == "marketdata" {
			if err := parseMarketDataToMap(line, resultMap); err != nil {
				return nil, err
			}

			continue
		}

		// header == "security"
		if err := parseSecurityToMap(line, resultMap); err != nil {
			return nil, err
		}
	}

	return filteredBondListFromMap(resultMap), nil
}

var (
	errSkip                      = errors.New("need skip this data line")
	errWrongMarketDataLineLength = errors.New("wrong market data line length")
)

func filteredBondListFromMap(resultMap map[string]Bond) []Bond {
	result := make([]Bond, 0, len(resultMap))
	emptySecPart := SecurityPart{}

	for _, bond := range resultMap {
		if bond.SecurityPart != emptySecPart {
			result = append(result, bond)
		}
	}

	return result
}

func parseMarketDataToMap(line []string, resultMap map[string]Bond) error {
	market, err := tryParseMarketData(line)

	if errors.Is(err, errSkip) {
		return nil
	}

	if err != nil {
		return err
	}

	if _, exist := resultMap[line[0]]; !exist {
		resultMap[line[0]] = Bond{
			ID:             line[0],
			MarketDataPart: market,
		}
	}

	return nil
}

func parseSecurityToMap(line []string, resultMap map[string]Bond) error {
	security, err := tryParseSecurity(line)

	if errors.Is(err, errSkip) {
		return nil
	}

	if err != nil {
		return err
	}

	if bond, exist := resultMap[line[0]]; exist {
		bond.SecurityPart = security

		resultMap[line[0]] = bond
	}

	return nil
}

func tryParseMarketData(line []string) (MarketDataPart, error) {
	const marketDataLinesCount = 2

	if len(line) != marketDataLinesCount {
		return MarketDataPart{}, fmt.Errorf("tryParseMarketData: %w: %d", errWrongMarketDataLineLength, len(line))
	}

	if line[1] == "" {
		return MarketDataPart{}, errSkip
	}

	currentPrice, err := strconv.ParseFloat(line[1], 64)
	if err != nil {
		return MarketDataPart{}, fmt.Errorf("cannot parse MarketData %w", err)
	}

	return MarketDataPart{currentPrice}, nil
}

var errWrongSecurityLineLength = errors.New("wrong security line length")

//nolint:funlen,cyclop
func tryParseSecurity(line []string) (SecurityPart, error) {
	const securityLinesCount = 12

	if len(line) != securityLinesCount {
		return SecurityPart{}, fmt.Errorf("tryParseSecurity: %w: %d", errWrongSecurityLineLength, len(line))
	}

	shortName := line[1]

	coupon, err := strconv.ParseFloat(line[2], 64)
	if err != nil {
		return SecurityPart{}, fmt.Errorf("cannot parse Security %w", err)
	}

	nextCoupon, err := util.ParseMoexDate(line[3])
	if err != nil {
		return SecurityPart{}, fmt.Errorf("cannot parse Security %w", err)
	}

	accCoupon, err := strconv.ParseFloat(line[4], 64)
	if err != nil {
		return SecurityPart{}, fmt.Errorf("cannot parse Security %w", err)
	}

	prevPrice, err := datastruct.ParseOptionalFloat64(line[5])
	if err != nil {
		return SecurityPart{}, fmt.Errorf("cannot parse Security %w", err)
	}

	prevPriceF, exist := prevPrice.Get()
	if !exist {
		return SecurityPart{}, errSkip
	}

	value, err := strconv.ParseFloat(line[6], 64)
	if err != nil {
		return SecurityPart{}, fmt.Errorf("cannot parse Security %w", err)
	}

	couponPeriod, err := strconv.ParseUint(line[7], 10, 64)
	if err != nil {
		return SecurityPart{}, fmt.Errorf("cannot parse Security %w", err)
	}

	priceStep, err := strconv.ParseFloat(line[8], 64)
	if err != nil {
		return SecurityPart{}, fmt.Errorf("cannot parse Security %w", err)
	}

	couponPercent, err := datastruct.ParseOptionalFloat64(line[9])
	if err != nil {
		return SecurityPart{}, fmt.Errorf("cannot parse Security %w", err)
	}

	endDate, err := util.ParseMoexDate(line[10])
	if err != nil {
		return SecurityPart{}, fmt.Errorf("cannot parse Security %w", err)
	}

	currency := line[11]
	if currency != "SUR" {
		return SecurityPart{}, errSkip
	}

	return SecurityPart{
		ShortName:        shortName,
		Coupon:           coupon,
		NextCoupon:       nextCoupon,
		AccCoupon:        accCoupon,
		PrevPricePercent: prevPriceF,
		Value:            value,
		CouponPeriod:     uint(couponPeriod),
		PriceStep:        priceStep,
		CouponPercent:    couponPercent,
		EndDate:          endDate,
		Currency:         currency,
	}, nil
}
