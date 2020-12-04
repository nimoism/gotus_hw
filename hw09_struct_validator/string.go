package hw09_struct_validator //nolint:golint,stylecheck

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type strValidate func(string, string, string) error

var strValidates = map[string]strValidate{
	tagLen:    strLenValidate,
	tagRegexp: strRegexpValidate,
	tagIn:     strInValidate,
}

type ErrLenNotEqual struct {
	Len      int
	Expected int
}

func (e ErrLenNotEqual) Error() string {
	return fmt.Sprintf("len %v != %v", e.Len, e.Expected)
}

var strLenValidate strValidate = func(name, constraint, value string) error {
	l, err := strconv.Atoi(constraint)
	if err != nil {
		return fmt.Errorf("str len format error: %w", err)
	}
	if vl := len(value); vl != l {
		return ValidationError{
			Field: name,
			Err:   ErrLenNotEqual{Len: vl, Expected: l},
		}
	}
	return nil
}

type ErrRegexpNotMatch struct {
	Value  string
	Regexp *regexp.Regexp
}

func (e ErrRegexpNotMatch) Error() string {
	return fmt.Sprintf("'%v' not matches '%v'", e.Value, e.Regexp.String())
}

var strRegexpCache = map[string]*regexp.Regexp{}

var strRegexpValidate strValidate = func(name, re, value string) error {
	compliedRegexp, ok := strRegexpCache[re]
	if !ok {
		compliedRegexp = regexp.MustCompile(re)
		strRegexpCache[re] = compliedRegexp
	}
	if !compliedRegexp.MatchString(value) {
		return ValidationError{
			Field: name,
			Err:   ErrRegexpNotMatch{Value: value, Regexp: compliedRegexp},
		}
	}
	return nil
}

type ErrStrNotIn struct {
	Value string
	Items []string
}

func (e ErrStrNotIn) Error() string {
	return fmt.Sprintf("%v not in %v", e.Value, e.Items)
}

var strInCache = map[string]map[string]struct{}{}

var strInValidate strValidate = func(name, constraint, value string) error {
	itemsMap, ok := strInCache[constraint]
	if !ok {
		itemList := strings.Split(constraint, tagInSep)
		itemsMap = make(map[string]struct{}, len(itemList))
		for _, item := range itemList {
			itemsMap[item] = struct{}{}
		}
		strInCache[constraint] = itemsMap
	}
	if _, ok := itemsMap[value]; !ok {
		items := make([]string, 0, len(itemsMap))
		for item := range itemsMap {
			items = append(items, item)
		}
		sort.Strings(items)
		return ValidationError{
			Field: name,
			Err:   ErrStrNotIn{Value: value, Items: items},
		}
	}
	return nil
}
