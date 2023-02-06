package main

import (
	"fmt"
	"log"
)

func part1(filename string) string {
	var iter LineIterator
	err := iter.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer iter.Close()

	var blueprint Blueprint
	var total int
	for iter.Next() {
		blueprint.Parse(iter.Value())
		blueprint.Optimize(24)
		total += blueprint.Quality()
		//fmt.Printf("%3d: %3d\n", blueprint.ID, blueprint.Quality())
	}
	return fmt.Sprint(total)
}

func part2(filename string) string {
	var iter LineIterator
	err := iter.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer iter.Close()

	var blueprint Blueprint
	var line int
	result := 1
	for iter.Next() {
		blueprint.Parse(iter.Value())
		blueprint.Optimize(32)
		result *= blueprint.MaxGeodes()
		//fmt.Printf("%3d: %3d\n", blueprint.ID, blueprint.MaxGeodes())
		line++
		if line >= 3 {
			break
		}
	}
	return fmt.Sprint(result)
}
