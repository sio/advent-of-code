package main

import (
	"testing"
)

const sample = "sample.txt"

func TestSample(t *testing.T) {
	tests := []struct {
		worker func(string) string
		input  string
		result string
	}{
		{worker: part1, input: sample, result: "13"},
		{worker: part2, input: sample, result: "1"},
		{worker: part1, input: "sample2.txt", result: "88"},
		{worker: part2, input: "sample2.txt", result: "36"},
	}
	for i, test := range tests {
		got := test.worker(test.input)
		expected := test.result
		if got != expected {
			t.Errorf("sample: part %d expected %q, got %q", i%2+1, expected, got)
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
