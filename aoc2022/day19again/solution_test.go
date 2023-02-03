package main

import (
	"testing"

	"strings"
)

func TestSample(t *testing.T) {
	tests := []struct {
		worker func(string) string
		input  string
		result string
	}{
		{part1, "sample.txt", "33"}, // input.txt: 1258 is too low
		{part2, "sample.txt", ""},
	}
	for i, test := range tests {
		got := strings.TrimSpace(test.worker(test.input))
		want := strings.TrimSpace(test.result)
		if got != want {
			t.Errorf("sample: part %d want %q, got %q", i%2+1, want, got)
			if strings.Contains(got, "\n") {
				for pos, char := range got {
					exp := []rune(want)[pos]
					if char != exp {
						t.Logf(
							"mismatching character #%d: got %c, want %c",
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
		part1("sample.txt")
	}
}

func BenchmarkPart2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		part2("sample.txt")
	}
}
