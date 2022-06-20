package moex

import (
	"bonds_calculator/internal/model/moex"
	asserts "github.com/stretchr/testify/assert"
	"io"
	"strings"
	"sync"
	"testing"
)

var (
	allBondsUrl       = "https://iss.moex.com/iss/engines/stock/markets/bonds/securities.csv?iss.meta=off&iss.only=marketdata,securities&securities.columns=SECID,SHORTNAME,COUPONVALUE,NEXTCOUPON,ACCRUEDINT,PREVPRICE,FACEVALUE,COUPONPERIOD,MINSTEP,COUPONPERCENT&marketdata.columns=SECID,LCURRENTPRICE"
	allBondizationUrl = "https://iss.moex.com/iss/securities/${bond}/bondization.csv?limit=unlimited&iss.meta=off&iss.only=amortizations,coupons&amortizations.columns=amortdate,value&coupons.columns=coupondate,value"
)

func TestBondNullFields(t *testing.T) {
	t.Parallel()

	assert := asserts.New(t)
	resp, err := client.Get(allBondsUrl)
	assert.NoError(err, "getting all bonds")

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	buf, err := io.ReadAll(resp.Body)
	assert.NoError(err, "reading body")

	bonds, err := moex.ParseBondsCp1251(buf)
	assert.NoError(err, "unmarshalling bonds")

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

	resp, _ := client.Get(allBondsUrl)
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	buf, _ := io.ReadAll(resp.Body)

	bonds, _ := moex.ParseBondsCp1251(buf)

	wg := sync.WaitGroup{}

	requestChan := make(chan string)

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			for id := range requestChan {
				resp, err := client.Get(strings.ReplaceAll(allBondizationUrl, "${bond}", id))
				assert.NoError(err, "getting all bondizations")

				if resp.Body != nil {
					defer resp.Body.Close()
				}

				buf, err := io.ReadAll(resp.Body)
				assert.NoError(err, "read bondization body")

				bondization, err := moex.ParseBondization(id, buf)
				assert.NoError(err, "unmarshalling bondization")

				for _, amortization := range bondization.Amortizations {
					assert.NotEqual(0, amortization.Value)
				}
			}

			wg.Done()
		}()
	}

	for _, bond := range bonds {
		requestChan <- bond.Id
	}

	close(requestChan)

	wg.Wait()
}
