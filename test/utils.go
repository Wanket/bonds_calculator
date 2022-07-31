package test

import (
	"bonds_calculator/internal/model/db"
	"errors"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	asserts "github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

var ErrTest = errors.New("test error")

func PrepareTest(t *testing.T) (*asserts.Assertions, *gomock.Controller) {
	t.Helper()

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

func isZero[T any](obj any) bool {
	if obj == nil {
		return true
	}

	switch reflect.TypeOf(obj).Kind() { //nolint:exhaustive
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(obj).IsNil()
	}

	var zero T

	return reflect.DeepEqual(obj, zero)
}

var (
	errHistoryUnsorted                   = errors.New("buy history is unsorted")
	errHistoryDateAfterEndDate           = errors.New("buy history date after end date")
	errHistoryCountZero                  = errors.New("buy history count is zero")
	errHistoryAccCouponNegative          = errors.New("buy history acc coupon is negative")
	errHistoryBondIDEmpty                = errors.New("buy history bond id is empty")
	errHistoryNominalValueZeroOrNegative = errors.New("buy history nominal value is zero or negative")
)

func CheckBuyHistoryValid(history []db.BuyHistory, endDate time.Time) error {
	if !slices.IsSortedFunc(history, func(left, right db.BuyHistory) bool {
		return left.Date.Before(right.Date)
	}) {
		return errHistoryUnsorted
	}

	for _, hst := range history {
		if !hst.Date.Before(endDate) {
			return errHistoryDateAfterEndDate
		}

		if hst.Count == 0 {
			return errHistoryCountZero
		}

		if hst.AccCoupon < 0 {
			return errHistoryAccCouponNegative
		}

		if hst.BondID == "" {
			return errHistoryBondIDEmpty
		}

		if hst.NominalValue <= 0 {
			return errHistoryNominalValueZeroOrNegative
		}
	}

	return nil
}
