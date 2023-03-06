package main

import (
	"testing"
)

func TestConversion(t *testing.T) {
	tests := map[string]int{
		// SNAFU: Decimal
		"1=-0-2": 1747,
		"12111":  906,
		"2=0=":   198,
		"21":     11,
		"2=01":   201,
		"111":    31,
		"20012":  1257,
		"112":    32,
		"1=-1=":  353,
		"1-12":   107,
		"12":     7,
		"1=":     3,
		"122":    37,

		"1":             1,
		"2":             2,
		"1-":            4,
		"10":            5,
		"11":            6,
		"2=":            8,
		"2-":            9,
		"20":            10,
		"1=0":           15,
		"1-0":           20,
		"1=11-2":        2022,
		"1-0---0":       12345,
		"1121-1110-1=0": 314159265,

		"2=-1=0": 4890,
	}
	for input, number := range tests {
		var got, want SnafuNumber
		err := got.Parse(input)
		if err != nil {
			t.Errorf("failed to parse %s: %v", input, err)
		}
		want = SnafuNumber(number)
		if got != want {
			t.Errorf("conversion to decimal failed for %q: got %d, want %d", input, got, want)
		}
		if want.String() != input {
			t.Errorf("conversion to SNAFU failed for %d: got %s, want %s", number, want.String(), input)
		}
	}
}
