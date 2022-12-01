package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

type ElfBag struct {
	OwnerID  int
	Items    int
	Calories int
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var line string
	var current, biggest ElfBag
	for scanner.Scan() {
		line = scanner.Text()
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
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	log.Printf(
		"Biggest bag: OwnerID=%d, Items=%d, Calories=%d",
		biggest.OwnerID,
		biggest.Items,
		biggest.Calories,
	)
}
