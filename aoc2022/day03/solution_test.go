package main

import (
	"testing"
)

func TestLetterScores(t *testing.T) {
	table := map[rune]int{
		'a': 1,
		'b': 2,
		'p': 16,
		'z': 26,
		'A': 26 + 1,
		'B': 26 + 2,
		'L': 38,
		'P': 42,
		'Z': 26 + 26,
	}
	var got int
	for k, expected := range table {
		got = LetterScore(k)
		if expected != got {
			t.Errorf("Incorrect score for %q: expected %d, got %d", k, expected, got)
		}
	}
}
