package test

import (
	asserts "github.com/stretchr/testify/assert"
	"reflect"
)

func CheckFailed[T any](assert *asserts.Assertions, t T, err error) {
	tIsNil := isZero[T](t)

	if tIsNil && err == nil {
		assert.Fail("bonds and err are both nil")
	}

	if !tIsNil && err != nil {
		assert.Fail("bonds and err are both not nil")
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
