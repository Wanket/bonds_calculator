package db

import (
	"time"
)

type BuyHistory struct {
	BondId       string
	Count        uint
	Date         time.Time
	Price        float64
	AccCoupon    float64
	NominalValue float64
}
