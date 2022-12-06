package main

import (
	"strconv"
)

type SlidingWindow struct {
	items []rune
	size  int
}

func (sw *SlidingWindow) Push(value rune) {
	sw.items = append(sw.items, value)
	if len(sw.items) > sw.size {
		sw.items = sw.items[len(sw.items)-sw.size:]
	}
}

func (sw *SlidingWindow) Contains(value rune) bool {
	for _, item := range sw.items {
		if item == value {
			return true
		}
	}
	return false
}

func (sw *SlidingWindow) Full() bool {
	return len(sw.items) == sw.size
}

func (sw *SlidingWindow) Unique() bool {
	itemset := make(map[rune]bool)
	for _, item := range sw.items {
		if itemset[item] {
			return false
		}
		itemset[item] = true
	}
	return len(itemset) == len(sw.items)
}

func LocateStartOfPacket(input <-chan rune) (position int) {
	return LocateEndOfMark(input, 4)
}

func LocateStartOfMessage(input <-chan rune) (position int) {
	return LocateEndOfMark(input, 14)
}

func LocateEndOfMark(input <-chan rune, markSize int) (position int) {
	window := SlidingWindow{size: markSize}
	for char := range input {
		position += 1
		window.Push(char)
		if window.Full() && window.Unique() {
			return position
		}
	}
	return -1
}

func part1(filename string) string {
	result := LocateStartOfPacket(ReadChars(filename))
	return strconv.Itoa(result)
}

func part2(filename string) string {
	result := LocateStartOfMessage(ReadChars(filename))
	return strconv.Itoa(result)
}
