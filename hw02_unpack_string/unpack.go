package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var result strings.Builder
	var prev rune
	for _, r := range str {
		repeat, err := strconv.Atoi(string(r))
		if err == nil {
			if prev == 0 {
				return "", ErrInvalidString
			}
			if repeat == 0 {
				prevString := strings.TrimSuffix(result.String(), string(prev))
				result.Reset()
				result.WriteString(prevString)

				continue
			}

			result.WriteString(strings.Repeat(string(prev), repeat-1))
			prev = 0
			continue
		}

		result.WriteRune(r)
		prev = r
	}

	return result.String(), nil
}
