//nolint:gomnd
package calculator

import (
	"bonds_calculator/internal/model/db"
	"bonds_calculator/internal/model/moex"
	testmoex "bonds_calculator/test/model/moex"
	testdataloader "github.com/peteole/testdata-loader"
)

func LoadBuyHistoryVariable() db.BuyHistory {
	return db.BuyHistory{
		BondID:       "SU24021RMFS6",
		Date:         testmoex.ParseDate("2022-06-23"),
		Price:        999.53,
		AccCoupon:    20.86,
		NominalValue: 1000,
		Count:        1,
	}
}

func LoadBuyHistory() db.BuyHistory {
	return db.BuyHistory{
		BondID:       "RU000A100TL1",
		Date:         testmoex.ParseDate("2022-06-21"),
		Price:        4998.2,
		AccCoupon:    38.26,
		NominalValue: 4900,
		Count:        1,
	}
}

func LoadMultiplyBuyHistory() []db.BuyHistory {
	return []db.BuyHistory{
		{
			BondID:       "RU000A100TL1",
			Date:         testmoex.ParseDate("2022-06-09"),
			Price:        4999.83,
			AccCoupon:    16.11,
			NominalValue: 4900,
			Count:        2,
		},
		{
			BondID:       "RU000A100TL1",
			Date:         testmoex.ParseDate("2022-05-26"),
			Price:        5203.946,
			AccCoupon:    51.68,
			NominalValue: 5240,
			Count:        3,
		},
	}
}

func LoadBondizationVariable() moex.Bondization {
	bonds, err := moex.ParseBondsCp1251(testdataloader.GetTestFile("test/data/moex/bond_variable.csv"))
	if err != nil {
		panic(err)
	}

	bondizations, err := moex.ParseBondization(
		bonds[0].ID,
		testdataloader.GetTestFile("test/data/moex/bondization_variable.csv"),
	)

	if err != nil {
		panic(err)
	}

	return bondizations
}
