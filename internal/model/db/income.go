package db

import (
	"time"
)

const (
	None   = 0
	Coupon = 1 << iota
	Maturity
	Amortization
	Sale
)

type IncomeType int

type Income struct {
	BondID     string
	IncomeType IncomeType
	Value      float64
	Date       time.Time
}
