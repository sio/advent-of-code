package main

import (
	"fmt"
)

func part1(filename string) string {
	basin := &BlizzardBasin{}
	basin.Load(filename)
	fmt.Println(basin)

	search := Search{basin: basin}
	return fmt.Sprint(search.ShortestPath())
}

func part2(filename string) string {
	return ""
}
