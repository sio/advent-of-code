package main

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
)

type ElfBag struct {
	OwnerID  int
	Items    int
	Calories int
}

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

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var line string
	var current ElfBag

	const topBagsCount = 3
	topBags := make([]ElfBag, topBagsCount)

	for scanner.Scan() {
		line = scanner.Text()
		if len(line) == 0 {
			if current.Calories > topBags[0].Calories {
				topBags = append(topBags, current)
				sort.Sort(byCalories(topBags))
				topBags = topBags[len(topBags)-topBagsCount:]
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
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
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
