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

func TestSamplePart2(t *testing.T) {
	basin := &BlizzardBasin{}
	basin.Load("sample.txt")
	search := Search{basin: basin}

	var commute int
	commute += search.ShortestPath(
		basin.entrance,
		basin.exit,
		commute,
	)
	want := 18
	if commute != want {
		t.Errorf("first trip: got %d, want %d", commute, want)
	}

	commute += search.ShortestPath(
		basin.exit,
		basin.entrance,
		commute,
	)
	want = 18 + 23
	if commute != want {
		t.Errorf("second trip: got %d, want %d", commute, want)
	}

	commute += search.ShortestPath(
		basin.entrance,
		basin.exit,
		commute,
	)
	want = 18 + 23 + 13
	if commute != want {
		t.Errorf("third trip: got %d, want %d", commute, want)
	}
}
