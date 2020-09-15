package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	if input == "" {
		return "", nil
	}
	inputRunes := []rune(input)
	if unicode.IsDigit(inputRunes[0]) {
		return "", ErrInvalidString
	}

	var builder strings.Builder
	var char rune
	var nextChar rune
	var isEscaped bool

	for i := range inputRunes[:len(inputRunes)-1] {
		char = inputRunes[i]
		nextChar = inputRunes[i+1]
		if !isEscaped && char == '\\' {
			if !unicode.IsDigit(nextChar) && nextChar != '\\' { // only digits and slashes escaping allowed
				return "", ErrInvalidString
			}
			isEscaped = true
			continue
		}
		if !isEscaped && unicode.IsDigit(char) {
			if unicode.IsDigit(nextChar) {
				return "", ErrInvalidString
			}
			continue
		}
		if unicode.IsDigit(nextChar) {
			repeatCount, err := strconv.Atoi(string(nextChar))
			if err != nil {
				return "", err
			}
			builder.WriteString(strings.Repeat(string(char), repeatCount))
		} else {
			builder.WriteString(string(char))
		}
		isEscaped = false
	}
	// handle last char
	if isEscaped || !unicode.IsDigit(nextChar) {
		builder.WriteString(string(nextChar))
	}

	return builder.String(), nil
}
