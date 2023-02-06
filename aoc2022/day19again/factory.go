package main

import (
	"fmt"
)

type Blueprint struct {
	ID            int
	Cost          [numberOfResources]Resources
	maxGeodeStock ResourceValue
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

func (b *Blueprint) Optimize(moves int) {
	// Initial factory state
	factory := FactoryState{
		blueprint: b,
		runway:    moves,
	}
	factory.output[Ore] = 1

	// Prepare result storage
	b.maxGeodeStock = 0

	// Start recursion
	factory.Optimize()
}

func (b *Blueprint) Quality() int {
	return b.ID * int(b.maxGeodeStock)
}

type FactoryState struct {
	blueprint *Blueprint
	runway    int
	output    Resources
	stock     Resources
}

// Build next robot
func (factory FactoryState) Produce(robot ResourceKind) FactoryState {

	cost := factory.blueprint.Cost[robot]

	const endgameMoves = 2 // one move to build the last robot + one move to benefit from it
	endgame := factory.stock.Add(factory.output.Times(factory.runway - endgameMoves))
	if factory.runway < endgameMoves || !endgame.Covers(cost) { // we can not build this robot factory
		factory.stock = factory.stock.Add(factory.output.Times(factory.runway))
		factory.runway = 0
		return factory
	}

	var done bool
	for factory.runway >= endgameMoves {
		factory.runway--
		produced := factory.output
		if factory.stock.Covers(cost) {
			factory.output[robot]++
			factory.stock = factory.stock.Sub(cost)
			done = true
		}
		factory.stock = factory.stock.Add(produced)
		if done {
			return factory
		}
	}
	panic("endgame estimation failed to short-circuit")
}

// Check if this search branch is exhausted
func (factory FactoryState) Done() bool {
	if factory.runway < 0 {
		panic("negative runway!")
	}
	return factory.runway == 0
}

// Find optimal blueprint output
func (factory FactoryState) Optimize() {
	if factory.stock[Geode] > factory.blueprint.maxGeodeStock {
		factory.blueprint.maxGeodeStock = factory.stock[Geode]
	}

	if factory.Done() {
		return // stop iteration
	}

	ceiling := factory.stock[Geode] + factory.output[Geode]*ResourceValue(factory.runway)
	for i := ResourceValue(factory.runway - 1); i > 0; i-- {
		ceiling += i
	}
	if ceiling < factory.blueprint.maxGeodeStock {
		return // short-circuit
	}

	for _, robot := range [...]ResourceKind{Geode, Obsidian, Clay, Ore} {
		next := factory.Produce(robot)
		next.Optimize()
	}
}
