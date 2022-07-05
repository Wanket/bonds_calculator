package model

import (
	"bonds_calculator/internal/model/db"
	"bonds_calculator/internal/model/moex"
	"fmt"
	"time"
)

const (
	Current  = iota
	Maturity = iota
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

func (calculator *IncomeCalculator) CalcPercent(buyHistory []db.BuyHistory, setting IncomeSetting) (float64, error) {
	percent := 0.0
	var sumCount uint = 0
	for _, buy := range buyHistory {
		oneHistoryPercent, err := calculator.CalcPercentForOneBuyHistory(buy, setting)
		if err != nil {
			return 0, err
		}

		percent += oneHistoryPercent * float64(buy.Count)

		sumCount += buy.Count
	}

	if sumCount == 0 {
		return 0, fmt.Errorf("no buy history with count > 0")
	}

	return percent / float64(sumCount), nil
}

func (calculator *IncomeCalculator) CalcPercentForOneBuyHistory(buyHistory db.BuyHistory, setting IncomeSetting) (float64, error) {
	currentBuyPrice := buyHistory.Price + buyHistory.AccCoupon

	percent := 0.0
	if setting == Maturity {
		percent += (buyHistory.NominalValue - currentBuyPrice) / currentBuyPrice
	}

	couponInx := getCouponIndexAfterBuyDate(buyHistory.Date, calculator.Coupons)

	calculator.recalculateCurrentCouponIfNeeded(couponInx, buyHistory)

	avgCoupon := -1.0
	amortizationInx := 0
	for ; couponInx < len(calculator.Coupons); couponInx++ {
		var amortizationChanged bool
		amortizationInx, amortizationChanged = calculator.shiftAmortization(amortizationInx, couponInx)

		coupon, exist := calculator.Coupons[couponInx].Value.Get()
		if !exist {
			if avgCoupon == -1.0 {
				avgCoupon = calcAvgCoupon(calculator.Coupons[:couponInx+1])
			}

			coupon = avgCoupon
		}

		percent += coupon / currentBuyPrice

		if amortizationChanged {
			currentBuyPrice -= calculator.Amortizations[amortizationInx].Value

			if currentBuyPrice <= 0 {
				return 0, fmt.Errorf("wrong amortizations sum")
			}
		}
	}

	return calcRelativePercent(percent, buyHistory.Date, calculator.Amortizations[amortizationInx].Date), nil
}

func (calculator *IncomeCalculator) recalculateCurrentCouponIfNeeded(couponInx int, buyHistory db.BuyHistory) {
	if _, exist := calculator.Coupons[couponInx].Value.Get(); !exist && 0 < couponInx && couponInx < len(calculator.Coupons) {
		startDay := calculator.Coupons[couponInx-1].Date
		endDay := calculator.Coupons[couponInx].Date

		if buyHistory.Date == startDay || buyHistory.Date == endDay {
			return
		}

		calculatedNextCoupon := buyHistory.AccCoupon /
			(buyHistory.Date.Sub(startDay).Hours() / 24) *
			(endDay.Sub(startDay).Hours() / 24)

		calculator.Coupons[couponInx].Value.Set(calculatedNextCoupon)
	}
}

func (calculator *IncomeCalculator) shiftAmortization(startIndex int, couponInx int) (index int, changed bool) {
	for index = startIndex; index < len(calculator.Amortizations)-1 && calculator.Amortizations[index].Date.Before(calculator.Coupons[couponInx].Date); index++ {
		changed = true
	}

	return
}

func calcAvgCoupon(coupons []moex.Coupon) float64 {
	var sum float64
	for _, coupon := range coupons {
		if val, exist := coupon.Value.Get(); exist {
			sum += val
		}
	}
	return sum / float64(len(coupons))
}

func getCouponIndexAfterBuyDate(buyDate time.Time, coupons []moex.Coupon) int {
	couponInx := 0
	for couponInx < len(coupons)-1 && !coupons[couponInx].Date.After(buyDate) {
		couponInx++
	}

	return couponInx
}

func calcRelativePercent(percent float64, startDate time.Time, endDate time.Time) float64 {
	return percent * 365 / (endDate.Sub(startDate).Hours() / 24)
}
