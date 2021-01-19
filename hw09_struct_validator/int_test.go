package hw09_struct_validator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntMinValidate(t *testing.T) {
	for _, tst := range []struct {
		value       int
		min         string
		expectedErr error
	}{
		{value: 20, min: "10"},
		{value: 10, min: "10"},
		{value: 9, min: "10", expectedErr: ErrLessThan{Value: 9, Min: 10}},
	} {
		err := intMinValidate("fn", tst.min, tst.value)
		var vErr ValidationError
		if tst.expectedErr != nil && errors.As(err, &vErr) {
			require.Equal(t, tst.expectedErr, vErr.Err)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestIntMaxValidate(t *testing.T) {
	for _, tst := range []struct {
		value       int
		max         string
		expectedErr error
	}{
		{value: 0, max: "10"},
		{value: 10, max: "10"},
		{value: 11, max: "10", expectedErr: ErrGreaterThan{Value: 11, Max: 10}},
	} {
		err := intMaxValidate("fn", tst.max, tst.value)
		var vErr ValidationError
		if tst.expectedErr != nil && errors.As(err, &vErr) {
			require.Equal(t, tst.expectedErr, vErr.Err)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestIntInValidate(t *testing.T) {
	for _, tst := range []struct {
		value       int
		items       string
		expectedErr error
	}{
		{value: 10, items: "10,20,30"},
		{value: 20, items: "10 , 20 , 30"},
		{value: 40, items: "10,20,30", expectedErr: ErrIntNotIn{Value: 40, Items: []int{10, 20, 30}}},
		{value: 0, items: "10,20,30", expectedErr: ErrIntNotIn{Value: 0, Items: []int{10, 20, 30}}},
	} {
		err := intInValidate("fn", "10,20,30", tst.value)
		var vErr ValidationError
		if tst.expectedErr != nil && errors.As(err, &vErr) {
			require.Equal(t, tst.expectedErr, vErr.Err)
		} else {
			require.NoError(t, err)
		}
	}
}
