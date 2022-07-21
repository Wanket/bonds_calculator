package test

import (
	"bonds_calculator/internal/model/db"
	"fmt"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	asserts "github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

func PrepareTest(t *testing.T) (*asserts.Assertions, *gomock.Controller) {
	log.SetOutput(ioutil.Discard)

	t.Parallel()

	assert := asserts.New(t)

	mockController := gomock.NewController(t)

	return assert, mockController
}

func CheckFailed[T any](assert *asserts.Assertions, t T, err error) {
	tIsNil := isZero[T](t)

	if tIsNil && err == nil {
		assert.Fail("value and err are both nil")
	}

	if !tIsNil && err != nil {
		assert.Fail("value and err are both not nil")
	}

	if err != nil && err.Error() == "" {
		assert.Fail("err is empty")
	}
}

func isZero[T any](t any) bool {
	if t == nil {
		return true
	}

	switch reflect.TypeOf(t).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(t).IsNil()
	}

	var zero T
	return reflect.DeepEqual(t, zero)
}

func CheckBuyHistoryValid(history []db.BuyHistory, endDate time.Time) error {
	if !slices.IsSortedFunc(history, func(left, right db.BuyHistory) bool {
		return left.Date.Before(right.Date)
	}) {
		return fmt.Errorf("buy history must be sorted")
	}

	for _, hst := range history {
		if !hst.Date.Before(endDate) {
			return fmt.Errorf("history date is equal or after end date")
		}

		if hst.Count == 0 {
			return fmt.Errorf("history count is zero")
		}

		if hst.AccCoupon < 0 {
			return fmt.Errorf("history acc coupon is negative")
		}

		if hst.BondId == "" {
			return fmt.Errorf("history bond id is empty")
		}

		if hst.NominalValue <= 0 {
			return fmt.Errorf("history nominal value is zero or negative")
		}
	}

	return nil
}
