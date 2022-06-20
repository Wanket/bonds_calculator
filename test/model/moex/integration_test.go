package moex

import (
	"bonds_calculator/internal/api"
	asserts "github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestBondNullFields(t *testing.T) {
	t.Parallel()

	assert := asserts.New(t)

	client := api.NewMoexClient(1)
	defer client.Close()

	bonds, err := client.GetBonds()
	assert.NoError(err, "getting bonds")

	for _, bond := range bonds {
		assert.NotEqual(0, bond.Coupon)
		assert.NotEqual(0, bond.AccCoupon)
		assert.NotEqual(0, bond.PrevPrice)
		assert.NotEqual(0, bond.Value)
		assert.NotEqual(0, bond.CouponPeriod)
		assert.NotEqual(0, bond.PriceStep)
		assert.NotEqual(0, bond.CurrentPrice)
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
		go func(id string) {
			bondization, err := client.GetBondization(id)
			assert.NoError(err, "getting bondization")

			for _, amortization := range bondization.Amortizations {
				assert.NotEqual(0, amortization.Value)
			}

			wg.Done()
		}(bond.Id)
	}

	wg.Wait()
}
