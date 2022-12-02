package main

import (
	"log"
	"sort"
	"strconv"
)

type byCalories []ElfBag

func (s byCalories) Len() int {
	return len(s)
}
func (s byCalories) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byCalories) Less(i, j int) bool {
	return s[i].Calories < s[j].Calories
}

const topBagsCount = 3

func part2(input string) {
	log.Println("Day 1 Part 2")
	var line string
	var current ElfBag

	topBags := make([]ElfBag, topBagsCount)

	for line = range ReadLines(input) {
		if len(line) == 0 {
			if current.Calories > topBags[0].Calories {
				topBags = AppendBag(topBags, current)
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
	topBags = AppendBag(topBags, current)
	var sumCalories int
	for _, bag := range topBags {
		log.Printf(
			"bag: OwnerID=%d, Items=%d, Calories=%d\n",
			bag.OwnerID,
			bag.Items,
			bag.Calories,
		)
		sumCalories += bag.Calories
	}
	log.Printf("sum: %d\n", sumCalories)
}

func AppendBag(bags []ElfBag, bag ElfBag) []ElfBag {
	bags = append(bags, bag)
	sort.Sort(byCalories(bags))
	return bags[len(bags)-topBagsCount:]
}
