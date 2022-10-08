package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var result strings.Builder

	var toPrint rune
	isEscaped := false
	for _, currentRune := range str {
		if isEscaped || (currentRune != '\\' && !unicode.IsDigit(currentRune)) {
			if toPrint != 0 {
				result.WriteRune(toPrint)
			}
			toPrint = currentRune
			isEscaped = false

			continue
		}

		if currentRune == '\\' {
			isEscaped = true
			continue
		}

		toRepeat, err := strconv.Atoi(string(currentRune))
		if err != nil || toPrint == 0 {
			return "", ErrInvalidString
		}

		result.WriteString(strings.Repeat(string(toPrint), toRepeat))
		toPrint = 0
		isEscaped = false
	}

	if toPrint != 0 {
		result.WriteRune(toPrint)
	}

	return result.String(), nil
}
