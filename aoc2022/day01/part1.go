package main

import (
	"log"
	"strconv"
)

type ElfBag struct {
	OwnerID  int
	Items    int
	Calories int
}

func part1(input string) {
	log.Println("Day 1 Part 1")
	var current, biggest ElfBag
	for line := range ReadLines(input) {
		if len(line) == 0 {
			if current.Calories > biggest.Calories {
				biggest = current
			}
			current = ElfBag{OwnerID: current.OwnerID + 1}
			continue
		}
		number, err := strconv.Atoi(line)
		if err != nil {
			log.Fatal(err)
		}
		current.Items += 1
		current.Calories += number
	}
	log.Printf(
		"Biggest bag: OwnerID=%d, Items=%d, Calories=%d",
		biggest.OwnerID,
		biggest.Items,
		biggest.Calories,
	)
}
