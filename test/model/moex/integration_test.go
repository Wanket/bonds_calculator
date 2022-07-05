package moex

import (
	"bonds_calculator/internal/api"
	asserts "github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestBondNullFields(t *testing.T) {
	t.Parallel()

	assert := asserts.New(t)

	client := api.NewMoexClient(1)
	defer client.Close()

	bonds, err := client.GetBonds()
	assert.NoError(err, "getting bonds")

	for _, bond := range bonds {
		assert.NoError(bond.IsValid(), "checking bond")
	}
}

func TestLoadAllBondizations(t *testing.T) {
	t.Parallel()

	assert := asserts.New(t)

	client := api.NewMoexClient(50)
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
