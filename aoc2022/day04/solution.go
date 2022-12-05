package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type SectionRange struct {
	Start int
	End   int
}

func (sr *SectionRange) Parse(input string) {
	boundaries := strings.Split(input, "-")
	if len(boundaries) != 2 {
		log.Fatalf("invalid range definition: %q", input)
	}
	var err error
	sr.Start, err = strconv.Atoi(boundaries[0])
	if err != nil {
		log.Fatalf("could not parse lower boundary: %v", err)
	}
	sr.End, err = strconv.Atoi(boundaries[1])
	if err != nil {
		log.Fatalf("could not parse upper boundary: %v", err)
	}
}

func (this *SectionRange) Contains(other *SectionRange) bool {
	return this.Start <= other.Start && this.End >= other.End
}

func (this *SectionRange) Overlaps(other *SectionRange) bool {
	return this.Contains(other) || (this.Start >= other.Start && this.Start <= other.End) || (this.End >= other.Start && this.End <= other.End)
}

func part1(filename string) {
	var first, second *SectionRange
	first, second = new(SectionRange), new(SectionRange)
	var answer int
	for line := range ReadLines(filename) {
		elves := strings.Split(line, ",")
		if len(elves) != 2 {
			log.Fatalf("invalid input line: %q", line)
		}
		first.Parse(elves[0])
		second.Parse(elves[1])
		if first.Contains(second) || second.Contains(first) {
			answer += 1
		}
	}
	fmt.Printf("Part 1 answer: %d\n", answer)
}

func part2(filename string) {
	var first, second *SectionRange
	first, second = new(SectionRange), new(SectionRange)
	var answer int
	for line := range ReadLines(filename) {
		elves := strings.Split(line, ",")
		if len(elves) != 2 {
			log.Fatalf("invalid input line: %q", line)
		}
		first.Parse(elves[0])
		second.Parse(elves[1])
		if first.Overlaps(second) {
			//log.Println(line)
			answer += 1
		}
	}
	fmt.Printf("Part 2 answer: %d\n", answer)
}
