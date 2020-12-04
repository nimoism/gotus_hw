package hw09_struct_validator //nolint:golint,stylecheck

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type intValidate func(string, string, int) error

var intValidates = map[string]intValidate{
	tagMin: intMinValidate,
	tagMax: intMaxValidate,
	tagIn:  intInValidate,
}

type ErrGreaterThan struct {
	Value int
	Max   int
}

func (e ErrGreaterThan) Error() string {
	return fmt.Sprintf("%v > %v", e.Value, e.Max)
}

var intMinValidate intValidate = func(name, rawConstraint string, value int) error {
	constraint, err := strconv.Atoi(rawConstraint)
	if err != nil {
		return fmt.Errorf("int min format error: %w", err)
	}
	if value < constraint {
		return ValidationError{
			Field: name,
			Err:   ErrLessThan{Value: value, Min: constraint},
		}
	}
	return nil
}

type ErrLessThan struct {
	Value int
	Min   int
}

func (e ErrLessThan) Error() string {
	return fmt.Sprintf("%v < %v", e.Value, e.Min)
}

var intMaxValidate intValidate = func(name, rawConstraint string, value int) error {
	constraint, err := strconv.Atoi(rawConstraint)
	if err != nil {
		return fmt.Errorf("int max format error: %w", err)
	}
	if value > constraint {
		return ValidationError{
			Field: name,
			Err:   ErrGreaterThan{Value: value, Max: constraint},
		}
	}
	return nil
}

type ErrIntNotIn struct {
	Value int
	Items []int
}

func (e ErrIntNotIn) Error() string {
	return fmt.Sprintf("%v not in %v", e.Value, e.Items)
}

var intInCache = map[string]map[int]struct{}{}

var intInValidate intValidate = func(name, rawConstraint string, value int) error {
	itemsMap, ok := intInCache[rawConstraint]
	if !ok {
		rawItems := strings.Split(rawConstraint, tagInSep)
		itemsMap = make(map[int]struct{}, len(rawItems))
		for _, rawItem := range rawItems {
			item, err := strconv.Atoi(rawItem)
			if err != nil {
				return fmt.Errorf("int in format error: %w", err)
			}
			itemsMap[item] = struct{}{}
			intInCache[rawConstraint] = itemsMap
		}
	}
	if _, ok := itemsMap[value]; !ok {
		items := make([]int, 0, len(itemsMap))
		for item := range itemsMap {
			items = append(items, item)
		}
		sort.Ints(items)
		return ValidationError{
			Field: name,
			Err:   ErrIntNotIn{Value: value, Items: items},
		}
	}
	return nil
}
