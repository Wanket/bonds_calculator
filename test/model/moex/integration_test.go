package moex_test

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
	assert.NoError(err)

	for _, bond := range bonds {
		assert.NoError(bond.IsValid())
	}
}

func TestLoadAllBondizations(t *testing.T) {
	assert, _ := test.PrepareTest(t)

	client := api.NewMoexClient(25)
	defer client.Close()

	bonds, err := client.GetBonds()
	assert.NoError(err)

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(bonds))

	for _, bond := range bonds {
		go func(id string, endDate time.Time) {
			bondization, err := client.GetBondization(id)
			assert.NoError(err)

			assert.NoError(bondization.IsValid(endDate))

			waitGroup.Done()
		}(bond.ID, bond.EndDate)
	}

	waitGroup.Wait()
}
