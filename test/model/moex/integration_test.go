package moex_test

import (
	"bonds_calculator/internal/api"
	"bonds_calculator/test"
	"testing"
)

func TestBondNullFields(t *testing.T) {
	assert, _ := test.PrepareTest(t)

	client := api.NewMoexClient(1)
	defer client.Close()

	bonds, err := client.GetBonds()
	assert.NoError(err)

	for _, bond := range bonds {
		assert.NoError(bond.IsValid())
	}
}

func TestLoadAllBondizations(t *testing.T) {
	assert, _ := test.PrepareTest(t)

	client := api.NewMoexClient(10)
	defer client.Close()

	bonds, err := client.GetBonds()
	assert.NoError(err)

	for bondizationResult := range client.GetBondizationsAsync(bonds) {
		assert.NoError(bondizationResult.Error)

		assert.NoError(bondizationResult.Bondization.IsValid(bondizationResult.Bond.EndDate))
	}
}
