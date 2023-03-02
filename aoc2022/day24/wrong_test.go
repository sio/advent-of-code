package main

import (
	"testing"

	"strconv"
)

func TestWrongValue(t *testing.T) {
	result, err := strconv.Atoi(part2("input.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if result >= 933 {
		t.Errorf("answer for input.txt is too high: %d", result)
	}
}
