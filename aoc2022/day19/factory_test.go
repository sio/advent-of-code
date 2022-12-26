package main

import (
	"testing"

	"fmt"
)

func TestFactoryParser(t *testing.T) {
	expected := Factory{
		ID: 1,
		Blueprint: map[ResourcePack]ResourcePack{
			ResourcePack{0, 0, 0, 0}: ResourcePack{0, 0, 0, 0},
			ResourcePack{1, 0, 0, 0}: ResourcePack{4, 0, 0, 0},
			ResourcePack{0, 1, 0, 0}: ResourcePack{2, 0, 0, 0},
			ResourcePack{0, 0, 1, 0}: ResourcePack{3, 14, 0, 0},
			ResourcePack{0, 0, 0, 1}: ResourcePack{2, 0, 7, 0},
		},
	}
	got := Factory{}
	err := got.Parse("Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 2 ore. Each obsidian robot costs 3 ore and 14 clay. Each geode robot costs 2 ore and 7 obsidian.")
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%v", got) != fmt.Sprintf("%v", expected) {
		t.Errorf("parsing failed: expected %v, got %v", expected, got)
	}
}
