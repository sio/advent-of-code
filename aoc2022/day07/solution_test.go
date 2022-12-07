package main

import (
	"testing"
)

const sample = "sample.txt"

func TestSample(t *testing.T) {
	results := []string{
		"95437",
		"24933642",
	}
	workers := []func(string) string{
		part1,
		part2,
	}
	if len(results) != len(workers) {
		t.Fatal("mismatch between number of worker functions and expected results")
	}
	for i := 0; i < len(results); i++ {
		got := workers[i](sample)
		expected := results[i]
		if got != expected {
			t.Errorf("sample: part %d expected %q, got %q", i+1, expected, got)
		}
	}
}
