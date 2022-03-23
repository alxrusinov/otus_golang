package hw02unpackstring

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

var regexpPattern = `^\d|\d{2,}`

func getRegexp(pattern string) (*regexp.Regexp, error) {
	reg, err := regexp.Compile(pattern)
	return reg, err
}

func checkIsInValidString(str string, match *regexp.Regexp) bool {
	return match.MatchString(str)
}

func checkStringForEmpty(str string) bool {
	return len(str) == 0
}

func Unpack(value string) (string, error) {
	var builder strings.Builder
	runeSlice := []rune(value)
	length := len(runeSlice)
	goNext := false
	matchString, err := getRegexp(regexpPattern)
	if err != nil {
		return "", err
	}

	if checkStringForEmpty(value) {
		return "", nil
	}

	if checkIsInValidString(value, matchString) {
		return "", ErrInvalidString
	}

	for i, v := range runeSlice {
		if goNext {
			goNext = false
			continue
		}

		if i == 0 && unicode.IsDigit(v) {
			return "", ErrInvalidString
		}

		if i == length-1 {
			builder.WriteString(string(v))
			continue
		}

		nextVal := runeSlice[i+1]
		if !unicode.IsDigit(nextVal) {
			builder.WriteString(string(v))
			continue
		}
		count, err := strconv.Atoi(string(nextVal))
		if err != nil {
			return "", ErrInvalidString
		}

		str := strings.Repeat(string(v), count)

		builder.WriteString(str)
		goNext = true
	}

	return builder.String(), nil
}
