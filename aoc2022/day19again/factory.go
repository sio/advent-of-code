package main

import (
	"fmt"
)

type Blueprint struct {
	ID   int
	Cost [numberOfResources]Resources
}

func (b *Blueprint) Parse(line string) {
	_, err := fmt.Sscanf(
		line,
		"Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian.",
		&b.ID,
		&b.Cost[Ore][Ore],
		&b.Cost[Clay][Ore],
		&b.Cost[Obsidian][Ore],
		&b.Cost[Obsidian][Clay],
		&b.Cost[Geode][Ore],
		&b.Cost[Geode][Obsidian],
	)
	if err != nil {
		panic(err)
	}
}

type FactoryState struct {
	blueprint *Blueprint
	runway    int
	output    Resources
	stock     Resources
}

func (current FactoryState) Produce(robot ResourceKind) (next FactoryState) {
	return current
}
