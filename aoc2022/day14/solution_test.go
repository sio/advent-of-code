package main

import (
	"testing"

	"strings"
)

const sample = "sample.txt"

func TestSample(t *testing.T) {
	tests := []struct {
		worker func(string) string
		input  string
		result string
	}{
		{worker: part1, input: sample, result: "24"},
		{worker: part2, input: sample, result: ""},
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
