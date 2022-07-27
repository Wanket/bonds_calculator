package calculator

import (
	"bonds_calculator/internal/model/db"
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/internal/util"
	"errors"
	"time"
)

const (
	Current = IncomeSetting(iota)
	Maturity
)

type IncomeSetting int

type IncomeCalculator struct {
	Amortizations []moex.Amortization
	Coupons       []moex.Coupon
}

func NewIncomeCalculator(bondization *moex.Bondization) IncomeCalculator {
	return IncomeCalculator{
		Amortizations: bondization.Amortizations,
		Coupons:       bondization.Coupons,
	}
}

var errNoBuyHistoryWithPositiveCount = errors.New("no buy history with count > 0")

func (calculator *IncomeCalculator) CalcPercent(buyHistory []db.BuyHistory, setting IncomeSetting) (float64, error) {
	var (
		percent       = 0.0
		sumCount uint = 0
	)

	for _, buy := range buyHistory {
		oneHistoryPercent, err := calculator.CalcPercentForOneBuyHistory(buy, setting)
		if err != nil {
			return 0, err
		}

		percent += oneHistoryPercent * float64(buy.Count)

		sumCount += buy.Count
	}

	if sumCount == 0 {
		return 0, errNoBuyHistoryWithPositiveCount
	}

	return percent / float64(sumCount), nil
}

var errWrongAmortizationsSum = errors.New("wrong amortizations sum")

func (calculator *IncomeCalculator) CalcPercentForOneBuyHistory(
	buyHistory db.BuyHistory,
	setting IncomeSetting,
) (float64, error) {
	currentBuyPrice := buyHistory.Price + buyHistory.AccCoupon

	percent := 0.0
	if setting == Maturity {
		percent += (buyHistory.NominalValue - currentBuyPrice) / currentBuyPrice
	}

	couponInx := getCouponIndexAfterBuyDate(buyHistory.Date, calculator.Coupons)
	amortizationInx := getAmortizationsAfterBuyDate(buyHistory.Date, calculator.Amortizations)

	calculator.recalculateCurrentCouponIfNeeded(couponInx, buyHistory)

	avgCoupon := -1.0
	accReturned := false

	for ; couponInx < len(calculator.Coupons); couponInx++ {
		coupon, exist := calculator.Coupons[couponInx].Value.Get()
		if !exist {
			if avgCoupon == -1.0 {
				avgCoupon = calcAvgCoupon(calculator.Coupons[:couponInx+1])
			}

			coupon = avgCoupon
		}

		if !accReturned {
			currentBuyPrice -= buyHistory.AccCoupon

			coupon -= buyHistory.AccCoupon

			accReturned = true
		}

		percent += coupon / currentBuyPrice

		amortizationsSum := calculator.calculateAmortizationSum(&amortizationInx, couponInx)

		offsetNominalPercent := amortizationsSum / buyHistory.NominalValue

		currentBuyPrice -= buyHistory.Price * offsetNominalPercent

		if currentBuyPrice <= -0.0001 {
			return 0, errWrongAmortizationsSum
		}
	}

	if amortizationInx == len(calculator.Amortizations) {
		amortizationInx--
	}

	return calcRelativePercent(percent, buyHistory.Date, calculator.Amortizations[amortizationInx].Date), nil
}

func (calculator *IncomeCalculator) calculateAmortizationSum(amortizationInx *int, couponInx int) float64 {
	amortizationsSum := 0.0

	for ; *amortizationInx < len(calculator.Amortizations); *amortizationInx++ {
		if calculator.Amortizations[*amortizationInx].Date.After(calculator.Coupons[couponInx].Date) {
			break
		}

		amortizationsSum += calculator.Amortizations[*amortizationInx].Value
	}

	return amortizationsSum
}

func (calculator *IncomeCalculator) recalculateCurrentCouponIfNeeded(couponInx int, buyHistory db.BuyHistory) {
	_, exist := calculator.Coupons[couponInx].Value.Get()
	if !exist && 0 < couponInx && couponInx < len(calculator.Coupons) {
		startDay := calculator.Coupons[couponInx-1].Date
		endDay := calculator.Coupons[couponInx].Date

		if buyHistory.Date == startDay || buyHistory.Date == endDay {
			return
		}

		calculatedNextCoupon := buyHistory.AccCoupon /
			(buyHistory.Date.Sub(startDay).Hours() / util.DayMultiplier) *
			(endDay.Sub(startDay).Hours() / util.DayMultiplier)

		calculator.Coupons[couponInx].Value.Set(calculatedNextCoupon)
	}
}

func calcAvgCoupon(coupons []moex.Coupon) float64 {
	sum := .0
	count := 0

	for _, coupon := range coupons {
		if val, exist := coupon.Value.Get(); exist {
			sum += val
			count++
		}
	}

	return sum / float64(count)
}

func getCouponIndexAfterBuyDate(buyDate time.Time, coupons []moex.Coupon) int {
	couponInx := 0
	for couponInx < len(coupons)-1 && !coupons[couponInx].Date.After(buyDate) {
		couponInx++
	}

	return couponInx
}

func getAmortizationsAfterBuyDate(buyDate time.Time, amortizations []moex.Amortization) int {
	amortizationInx := 0
	for amortizationInx < len(amortizations)-1 && !amortizations[amortizationInx].Date.After(buyDate) {
		amortizationInx++
	}

	return amortizationInx
}

func calcRelativePercent(percent float64, startDate time.Time, endDate time.Time) float64 {
	return percent * 365 / (endDate.Sub(startDate).Hours() / util.DayMultiplier)
}
