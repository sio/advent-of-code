package main

import (
	"testing"

	"strconv"
	"strings"
)

const sample = "sample.txt"

func TestSample(t *testing.T) {
	tests := []struct {
		worker func(string) string
		input  string
		result string
	}{
		{worker: part1, input: sample, result: "13"},
		{worker: part2, input: sample, result: ""},
		{worker: part1, input: "sample2.txt", result: strconv.Itoa(0 + 3)},
	}
	for i, test := range tests {
		got := strings.TrimSpace(test.worker(test.input))
		expected := strings.TrimSpace(test.result)
		if got != expected {
			t.Errorf("sample: part %d expected %q, got %q", i%2+1, expected, got)
			if strings.Contains(got, "\n") {
				for pos, char := range got {
					exp := []rune(expected)[pos]
					if char != exp {
						t.Logf(
							"mismatching character #%d: got %c, expected %c",
							pos,
							char,
							exp,
						)
					}
				}
			}
		}
	}
}

func TestParsing(t *testing.T) {
	var files = []string{
		"sample.txt",
		"input.txt",
	}
	for _, filename := range files {
		for line := range ReadLines(filename) {
			if len(line) == 0 {
				continue
			}
			list := &NestedList{}
			err := list.Parse(line)
			if err != nil || list.String() != line {
				t.Errorf("parsing error (%v): %s -> %s", err, line, list.String())
			}
		}
	}
}

func TestGuessing(t *testing.T) {
	value, err := strconv.Atoi(part1("input.txt"))
	if err != nil {
		t.Errorf("number parsing error: %v", err)
	}
	if value >= 5788 {
		t.Errorf("part 1 guess is too high: %d", value)
	}
}

func BenchmarkPart1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		part1(sample)
	}
}

func BenchmarkPart2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		part2(sample)
	}
}
