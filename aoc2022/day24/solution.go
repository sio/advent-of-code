package main

import (
	"fmt"
)

func part1(filename string) string {
	basin := &BlizzardBasin{}
	basin.Load(filename)
	fmt.Println(basin)

	search := Search{basin: basin}
	commute := search.ShortestPath(
		basin.entrance,
		basin.exit,
		0,
	)
	return fmt.Sprint(commute)
}

func part2(filename string) string {
	basin := &BlizzardBasin{}
	basin.Load(filename)
	search := Search{basin: basin}

	var commute int
	commute += search.ShortestPath(
		basin.entrance,
		basin.exit,
		commute,
	)
	commute += search.ShortestPath(
		basin.exit,
		basin.entrance,
		commute,
	)
	commute += search.ShortestPath(
		basin.entrance,
		basin.exit,
		commute,
	)
	return fmt.Sprint(commute)
}
