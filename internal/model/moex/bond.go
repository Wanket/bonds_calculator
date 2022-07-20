package moex

import (
	"bonds_calculator/internal/model/datastuct"
	"bonds_calculator/internal/util"
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"io"
	"strconv"
	"time"
)

// Structs for this response:
// https://iss.moex.com/iss/engines/stock/markets/bonds/securities.csv?iss.meta=off&iss.only=marketdata,securities&securities.columns=$1&marketdata.columns=$2
// $1 = SECID,SHORTNAME,COUPONVALUE,NEXTCOUPON,ACCRUEDINT,PREVPRICE,FACEVALUE,COUPONPERIOD,MINSTEP,COUPONPERCENT,MATDATE
// $2 = SECID,LCURRENTPRICE

type (
	Bond struct {
		Id string
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
		CouponPercent    datastuct.Optional[float64]
		EndDate          time.Time
	}

	MarketDataPart struct {
		CurrentPricePercent float64
	}
)

func (bond *Bond) IsValid() error {
	if bond.Id == "" {
		return fmt.Errorf("bond id is empty")
	}

	if bond.ShortName == "" {
		return fmt.Errorf("bond short name is empty")
	}

	if bond.Coupon < 0 {
		return fmt.Errorf("bond coupon is <= 0")
	}

	if bond.AccCoupon < 0 {
		return fmt.Errorf("bond acc coupon is < 0")
	}

	if bond.Value <= 0 {
		return fmt.Errorf("bond value is <= 0")
	}

	if zeroCouponRepiod, zeroNextCoupon := bond.CouponPeriod == 0, bond.NextCoupon.IsZero(); zeroCouponRepiod != zeroNextCoupon {
		return fmt.Errorf("bond coupon period is 0 xor next coupon is zero")
	}

	if bond.PriceStep <= 0 {
		return fmt.Errorf("bond price step is <= 0")
	}

	return nil
}

func (bond *Bond) AbsoluteCurrentPrice() float64 {
	return bond.CurrentPricePercent * bond.Value
}

func ParseBondsCp1251(buf []byte) ([]Bond, error) {
	decoded, err := charmap.Windows1251.NewDecoder().Bytes(buf)
	if err != nil {
		return nil, err
	}

	return ParseBonds(decoded)
}

func ParseBonds(buf []byte) ([]Bond, error) {
	reader := NewReader(bytes.NewReader(buf))

	header := ""

	resultMap := make(map[string]Bond)

	for {
		line, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if len(line) == 1 {
			if line[0] != "marketdata" && line[0] != "securities" {
				return nil, fmt.Errorf("invalid header %s", line[0])
			}

			header = line[0]

			continue
		}

		if header == "marketdata" {
			market, err := tryParseMarketData(line)

			if errors.Is(err, skipError) {
				continue
			}

			if err != nil {
				return nil, err
			}

			if _, exist := resultMap[line[0]]; !exist {
				resultMap[line[0]] = Bond{
					Id:             line[0],
					MarketDataPart: market,
				}
			}

			continue
		}

		// header == "security"
		security, err := tryParseSecurity(line)

		if errors.Is(err, skipError) {
			continue
		}

		if err != nil {
			return nil, err
		}

		if bond, exist := resultMap[line[0]]; exist {
			bond.SecurityPart = security

			resultMap[line[0]] = bond
		}
	}

	var emptySecPart SecurityPart
	result := make([]Bond, 0, len(resultMap))
	for _, bond := range resultMap {
		if bond.SecurityPart != emptySecPart {
			result = append(result, bond)
		}
	}

	return result, nil
}

var skipError = errors.New("need skip this data line")

func tryParseMarketData(line []string) (MarketDataPart, error) {
	if len(line) != 2 {
		return MarketDataPart{}, fmt.Errorf("wrong MarketData data line len %d", len(line))
	}

	if line[1] == "" {
		return MarketDataPart{}, skipError
	}

	currentPrice, err := strconv.ParseFloat(line[1], 64)
	if err != nil {
		return MarketDataPart{}, fmt.Errorf("cannot parse MarketData %v", err)
	}

	return MarketDataPart{currentPrice}, nil
}

func tryParseSecurity(line []string) (SecurityPart, error) {
	if len(line) != 11 {
		return SecurityPart{}, fmt.Errorf("wrong Security data line len %d", len(line))
	}

	prevPrice, err := util.ParseOptionalFloat64(line[5])
	if _, exist := prevPrice.Get(); !exist && err == nil {
		return SecurityPart{}, skipError
	}

	endDate, err := time.Parse("2006-01-02", line[10])
	if err != nil && line[10] == "0000-00-00" {
		endDate = time.Time{}
	}

	shortName := line[1]
	coupon, err := strconv.ParseFloat(line[2], 64)
	nextCoupon, err := time.Parse("2006-01-02", line[3])
	accCoupon, err := strconv.ParseFloat(line[4], 64)
	prevPriceF, _ := prevPrice.Get()
	value, err := strconv.ParseFloat(line[6], 64)
	couponPeriod, err := strconv.ParseInt(line[7], 10, 64)
	priceStep, err := strconv.ParseFloat(line[8], 64)
	couponPercent, err := util.ParseOptionalFloat64(line[9])

	if err != nil {
		return SecurityPart{}, fmt.Errorf("cannot parse Security %v", err)
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
	}, nil
}
