package main

import (
	"fmt"
)

func part1(filename string) string {
	basin := &BlizzardBasin{}
	basin.Load(filename)
	fmt.Println(basin)
	for i := 0; i < 18; i++ {
		fmt.Println("Round", i)
		fmt.Println(basin.Render(i))
	}
	return ""
}

func part2(filename string) string {
	return ""
}
