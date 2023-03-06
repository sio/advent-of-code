package main

import (
	"fmt"
	"strings"
)

var snafuDigit = [...]rune{'=', '-', '0', '1', '2'}

const snafuOffset int = 2
const snafuBase = len(snafuDigit)

type SnafuNumber int

func (num SnafuNumber) String() string {
	var digit []rune
	var value = int(num)
	for value != 0 {
		value += snafuOffset
		digit = append(digit, snafuDigit[value%snafuBase])
		value = value / snafuBase
	}
	var builder strings.Builder
	for i := len(digit) - 1; i >= 0; i-- {
		builder.WriteRune(digit[i])
	}
	return builder.String()
}

func (num *SnafuNumber) Parse(input string) error {
	digit := make([]int, 0, len(input))
	var pos int
	for _, char := range input {
		pos++
		value, err := parseSnafuDigit(char)
		if err != nil {
			return fmt.Errorf("SNAFU number %q parser failed at position %d: %w", input, pos, err)
		}
		digit = append(digit, value)
	}

	var value int
	for pos, dig := range digit {
		value += dig * power(snafuBase, len(digit)-pos-1)
	}
	*num = SnafuNumber(value)
	return nil
}

func parseSnafuDigit(char rune) (digit int, err error) {
	for value, digit := range snafuDigit {
		if digit == char {
			return value - snafuOffset, nil
		}
	}
	return 0, fmt.Errorf("invalid digit: %c", char)
}

func power(base, exp int) int {
	var result int = 1
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}
