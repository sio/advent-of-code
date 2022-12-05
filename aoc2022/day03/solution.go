package main

import (
	"fmt"
	"log"
	"strings"
)

func LetterScore(r rune) int {
	const a = int('a')
	const A = int('A')
	const alphabet = 26

	lowercase := int(r) - a + 1
	if lowercase <= alphabet && lowercase > 0 {
		return lowercase
	}
	uppercase := int(r) - A + 1
	if uppercase <= alphabet && uppercase > 0 {
		return uppercase + alphabet
	}
	panic(fmt.Sprintf("unsupported character: %q (ascii=%d, upper=%d, lower=%d)", r, int(r), uppercase, lowercase))
	return 0
}

func part1(filename string) {
	var total int
	for line := range ReadLines(filename) {
		if len(line)%2 != 0 {
			log.Fatalf("odd number of items in rucksack %q: %d", line, len(line))
		}
		for _, r := range line[len(line)/2:] {
			if strings.ContainsRune(line[:len(line)/2], r) {
				total += LetterScore(r)
				break
			}
		}
	}
	fmt.Printf("Part 1 score: %d\n", total)
}

func part2(filename string) {}
