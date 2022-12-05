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

func part2(filename string) {
	group := make([]string, 3)
	var total, index int
	for line := range ReadLines(filename) {
		group[index%len(group)] = line
		index += 1
		if index%len(group) == 0 {
			total += GroupScore(group)
		}
	}
	fmt.Printf("Part 2 score: %d\n", total)
}

func GroupScore(group []string) int {
	for _, r := range group[0] {
		if strings.ContainsRune(group[1], r) && strings.ContainsRune(group[2], r) {
			return LetterScore(r)
		}
	}
	panic(fmt.Sprintf("could not find a badge for group %q", group))
}
