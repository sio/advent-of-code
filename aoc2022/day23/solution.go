package main

import (
	"fmt"
)

func part1(filename string) string {
	var elves = &ElfGroup{}
	elves.Load(filename)
	fmt.Println(elves)
	return ""
}

func part2(filename string) string {
	return ""
}
