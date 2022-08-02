package moex_test

import (
	"bonds_calculator/internal/api"
	"bonds_calculator/test"
	"os"
	"testing"
)

func TestBondNullFields(t *testing.T) {
	if _, exist := os.LookupEnv("CI"); exist {
		t.Skip("Moex integration tests is unstable in GitHub Actions")
	}

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
	if _, exist := os.LookupEnv("CI"); exist {
		t.Skip("Moex integration tests is unstable in GitHub Actions")
	}

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
