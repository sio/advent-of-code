package main

import (
	"fmt"
)

func part1(filename string) string {
	var elves = &ElfGroup{}
	elves.Load(filename)
	elves.Play(10)
	return fmt.Sprint(elves.Result())
}

func part2(filename string) string {
	var elves = &ElfGroup{}
	elves.Load(filename)

	const maxRounds = 10000
	result := elves.Play(maxRounds)
	if result == maxRounds {
		panic(fmt.Sprintf("movements did not cease after %d rounds", maxRounds))
	}
	return fmt.Sprint(result + 1)
}
