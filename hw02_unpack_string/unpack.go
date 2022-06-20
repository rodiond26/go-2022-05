package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if str == "" {
		return "", nil
	}

	runes := []rune(str)

	if isDigit(runes[0]) {
		return "", ErrInvalidString
	}
	if len(runes) == 1 && isNotDigit(runes[0]) {
		return str, nil
	}

	var b strings.Builder

	for i := 0; i < len(runes)-1; i++ {
		if isDigit(runes[i]) && isDigit(runes[i+1]) {
			return "", ErrInvalidString
		}
		if isNotDigit(runes[i]) {
			if isNotDigit(runes[i+1]) {
				b.WriteString(string(runes[i]))
				continue
			} else {
				b.WriteString(repeate((runes[i]), (runes[i+1])))
			}
		}
	}

	lastRune := runes[len(runes)-1]
	if isNotDigit(lastRune) {
		b.WriteString(string(lastRune))
	}

	return b.String(), nil
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func isNotDigit(r rune) bool {
	return !unicode.IsDigit(r)
}

func repeate(r rune, count rune) string {
	return strings.Repeat(string(r), int(count-'0'))
}
