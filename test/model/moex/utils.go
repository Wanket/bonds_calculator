package moex

import (
	"bonds_calculator/internal/model/datastruct"
	"bonds_calculator/internal/model/moex"
	"time"
)

func LoadParsedBondization() moex.Bondization {
	return moex.Bondization{
		ID:            "RU000A100TL1",
		Amortizations: loadAmortizations(),
		Coupons:       loadCoupons(),
	}
}

func ParseDate(str string) time.Time {
	res, err := time.Parse("2006-01-02", str)
	if err != nil {
		panic(err)
	}

	return res
}

//nolint:gomnd
func loadCoupons() []moex.Coupon {
	return []moex.Coupon{
		{Date: ParseDate("2019-10-16"), Value: datastruct.NewOptional(123.29)},
		{Date: ParseDate("2019-11-15"), Value: datastruct.NewOptional(123.29)},
		{Date: ParseDate("2019-12-15"), Value: datastruct.NewOptional(123.29)},
		{Date: ParseDate("2020-01-14"), Value: datastruct.NewOptional(123.29)},
		{Date: ParseDate("2020-02-13"), Value: datastruct.NewOptional(123.29)},
		{Date: ParseDate("2020-03-14"), Value: datastruct.NewOptional(123.29)},
		{Date: ParseDate("2020-04-13"), Value: datastruct.NewOptional(123.29)},
		{Date: ParseDate("2020-05-13"), Value: datastruct.NewOptional(123.29)},
		{Date: ParseDate("2020-06-12"), Value: datastruct.NewOptional(123.29)},
		{Date: ParseDate("2020-07-12"), Value: datastruct.NewOptional(123.29)},
		{Date: ParseDate("2020-08-11"), Value: datastruct.NewOptional(123.29)},
		{Date: ParseDate("2020-09-10"), Value: datastruct.NewOptional(123.29)},
		{Date: ParseDate("2020-10-10"), Value: datastruct.NewOptional(123.29)},
		{Date: ParseDate("2020-11-09"), Value: datastruct.NewOptional(123.29)},
		{Date: ParseDate("2020-12-09"), Value: datastruct.NewOptional(123.29)},
		{Date: ParseDate("2021-01-08"), Value: datastruct.NewOptional(123.29)},
		{Date: ParseDate("2021-02-07"), Value: datastruct.NewOptional(123.29)},
		{Date: ParseDate("2021-03-09"), Value: datastruct.NewOptional(123.29)},
		{Date: ParseDate("2021-04-08"), Value: datastruct.NewOptional(123.29)},
		{Date: ParseDate("2021-05-08"), Value: datastruct.NewOptional(119.1)},
		{Date: ParseDate("2021-06-07"), Value: datastruct.NewOptional(114.9)},
		{Date: ParseDate("2021-07-07"), Value: datastruct.NewOptional(110.71)},
		{Date: ParseDate("2021-08-06"), Value: datastruct.NewOptional(106.52)},
		{Date: ParseDate("2021-09-05"), Value: datastruct.NewOptional(102.33)},
		{Date: ParseDate("2021-10-05"), Value: datastruct.NewOptional(98.14)},
		{Date: ParseDate("2021-11-04"), Value: datastruct.NewOptional(93.95)},
		{Date: ParseDate("2021-12-04"), Value: datastruct.NewOptional(89.75)},
		{Date: ParseDate("2022-01-03"), Value: datastruct.NewOptional(85.56)},
		{Date: ParseDate("2022-02-02"), Value: datastruct.NewOptional(81.37)},
		{Date: ParseDate("2022-03-04"), Value: datastruct.NewOptional(77.18)},
		{Date: ParseDate("2022-04-03"), Value: datastruct.NewOptional(72.99)},
		{Date: ParseDate("2022-05-03"), Value: datastruct.NewOptional(68.79)},
		{Date: ParseDate("2022-06-02"), Value: datastruct.NewOptional(64.6)},
		{Date: ParseDate("2022-07-02"), Value: datastruct.NewOptional(60.41)},
		{Date: ParseDate("2022-08-01"), Value: datastruct.NewOptional(56.22)},
		{Date: ParseDate("2022-08-31"), Value: datastruct.NewOptional(52.03)},
		{Date: ParseDate("2022-09-30"), Value: datastruct.NewOptional(47.84)},
		{Date: ParseDate("2022-10-30"), Value: datastruct.NewOptional(43.64)},
		{Date: ParseDate("2022-11-29"), Value: datastruct.NewOptional(39.45)},
		{Date: ParseDate("2022-12-29"), Value: datastruct.NewOptional(35.26)},
		{Date: ParseDate("2023-01-28"), Value: datastruct.NewOptional(31.07)},
		{Date: ParseDate("2023-02-27"), Value: datastruct.NewOptional(26.88)},
		{Date: ParseDate("2023-03-29"), Value: datastruct.NewOptional(22.68)},
		{Date: ParseDate("2023-04-28"), Value: datastruct.NewOptional(18.49)},
		{Date: ParseDate("2023-05-28"), Value: datastruct.NewOptional(14.3)},
		{Date: ParseDate("2023-06-27"), Value: datastruct.NewOptional(10.11)},
		{Date: ParseDate("2023-07-27"), Value: datastruct.NewOptional(5.92)},
		{Date: ParseDate("2023-08-26"), Value: datastruct.NewOptional(1.73)},
	}
}

//nolint:gomnd
func loadAmortizations() []moex.Amortization {
	return []moex.Amortization{
		{Date: ParseDate("2021-04-08"), Value: 340},
		{Date: ParseDate("2021-05-08"), Value: 340},
		{Date: ParseDate("2021-06-07"), Value: 340},
		{Date: ParseDate("2021-07-07"), Value: 340},
		{Date: ParseDate("2021-08-06"), Value: 340},
		{Date: ParseDate("2021-09-05"), Value: 340},
		{Date: ParseDate("2021-10-05"), Value: 340},
		{Date: ParseDate("2021-11-04"), Value: 340},
		{Date: ParseDate("2021-12-04"), Value: 340},
		{Date: ParseDate("2022-01-03"), Value: 340},
		{Date: ParseDate("2022-02-02"), Value: 340},
		{Date: ParseDate("2022-03-04"), Value: 340},
		{Date: ParseDate("2022-04-03"), Value: 340},
		{Date: ParseDate("2022-05-03"), Value: 340},
		{Date: ParseDate("2022-06-02"), Value: 340},
		{Date: ParseDate("2022-07-02"), Value: 340},
		{Date: ParseDate("2022-08-01"), Value: 340},
		{Date: ParseDate("2022-08-31"), Value: 340},
		{Date: ParseDate("2022-09-30"), Value: 340},
		{Date: ParseDate("2022-10-30"), Value: 340},
		{Date: ParseDate("2022-11-29"), Value: 340},
		{Date: ParseDate("2022-12-29"), Value: 340},
		{Date: ParseDate("2023-01-28"), Value: 340},
		{Date: ParseDate("2023-02-27"), Value: 340},
		{Date: ParseDate("2023-03-29"), Value: 340},
		{Date: ParseDate("2023-04-28"), Value: 340},
		{Date: ParseDate("2023-05-28"), Value: 340},
		{Date: ParseDate("2023-06-27"), Value: 340},
		{Date: ParseDate("2023-07-27"), Value: 340},
		{Date: ParseDate("2023-08-26"), Value: 140},
	}
}
