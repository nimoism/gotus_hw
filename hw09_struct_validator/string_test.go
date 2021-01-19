package hw09_struct_validator

import (
	"errors"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringLenValidate(t *testing.T) {
	for _, tst := range []struct {
		value       string
		len         string
		expectedErr error
	}{
		{value: "", len: "5", expectedErr: ErrLenNotEqual{Len: 0, Expected: 5}},
		{value: "1234", len: "5", expectedErr: ErrLenNotEqual{Len: 4, Expected: 5}},
		{value: "12345", len: "5", expectedErr: nil},
		{value: "123456", len: "5", expectedErr: ErrLenNotEqual{Len: 6, Expected: 5}},
	} {
		err := strLenValidate("fn", tst.len, tst.value)
		var vErr ValidationError
		if tst.expectedErr != nil && errors.As(err, &vErr) {
			require.Equal(t, tst.expectedErr, vErr.Err)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestStringRegexpValidate(t *testing.T) {
	re := `^\d{3}\w{2,4}$`
	compiledRe := regexp.MustCompile(re)
	for _, tst := range []struct {
		value       string
		re          string
		expectedErr error
	}{
		{value: "", re: re, expectedErr: ErrRegexpNotMatch{Value: "", Regexp: compiledRe}},
		{value: "12a", re: re, expectedErr: ErrRegexpNotMatch{Value: "12a", Regexp: compiledRe}},
		{value: "123abc", re: re, expectedErr: nil},
		{value: "1234abcde", re: re, expectedErr: ErrRegexpNotMatch{Value: "1234abcde", Regexp: compiledRe}},
	} {
		err := strRegexpValidate("fn", tst.re, tst.value)
		var vErr ValidationError
		if tst.expectedErr != nil && errors.As(err, &vErr) {
			require.Equal(t, tst.expectedErr, vErr.Err)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestStringInValidate(t *testing.T) {
	for _, tst := range []struct {
		value       string
		items       string
		expectedErr error
	}{
		{value: "foo", items: "foo,bar,baz"},
		{value: "bar", items: "foo , bar , baz"},
		{value: "wrong", items: "foo,bar,baz", expectedErr: ErrStrNotIn{Value: "wrong", Items: []string{"bar", "baz", "foo"}}},
		{value: "", items: "foo,bar,baz", expectedErr: ErrStrNotIn{Value: "", Items: []string{"bar", "baz", "foo"}}},
	} {
		err := strInValidate("fn", "foo,bar,baz", tst.value)
		var vErr ValidationError
		if tst.expectedErr != nil && errors.As(err, &vErr) {
			require.Equal(t, tst.expectedErr, vErr.Err)
		} else {
			require.NoError(t, err)
		}
	}
}
