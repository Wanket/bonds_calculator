package db

import (
	"time"
)

type BuyHistory struct {
	BondID       string
	Count        uint
	Date         time.Time
	Price        float64
	AccCoupon    float64
	NominalValue float64
}
