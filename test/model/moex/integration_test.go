package moex

import (
	"bonds_calculator/internal/api"
	"bonds_calculator/test"
	"sync"
	"testing"
	"time"
)

func TestBondNullFields(t *testing.T) {
	assert, _ := test.PrepareTest(t)

	client := api.NewMoexClient(1)
	defer client.Close()

	bonds, err := client.GetBonds()
	assert.NoError(err, "getting bonds")

	for _, bond := range bonds {
		assert.NoError(bond.IsValid(), "checking bond")
	}
}

func TestLoadAllBondizations(t *testing.T) {
	assert, _ := test.PrepareTest(t)

	client := api.NewMoexClient(25)
	defer client.Close()

	bonds, _ := client.GetBonds()

	wg := sync.WaitGroup{}
	wg.Add(len(bonds))

	for _, bond := range bonds {
		go func(id string, endDate time.Time) {
			bondization, err := client.GetBondization(id)
			assert.NoError(err, "getting bondization")

			assert.NoError(bondization.IsValid(endDate), "checking bondization")

			wg.Done()
		}(bond.Id, bond.EndDate)
	}

	wg.Wait()
}
