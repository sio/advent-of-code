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
	return ""
}
