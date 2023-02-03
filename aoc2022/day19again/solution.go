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
		fmt.Println(blueprint)
		blueprint.Optimize(24)
		fmt.Println(blueprint.Quality())
		total += blueprint.Quality()
	}
	return fmt.Sprint(total)
}

func part2(filename string) string {
	return ""
}
