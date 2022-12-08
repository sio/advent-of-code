package main

import (
	"log"
	"strconv"
)

type TreeHeight uint8

type Location struct {
	X uint
	Y uint
}

type Map struct {
	maxX uint
	maxY uint
	grid map[Location]TreeHeight
}

func (m *Map) Get(location Location) TreeHeight {
	return m.grid[location]
}

func (m *Map) Set(location Location, value TreeHeight) {
	if location.X > m.maxX {
		m.maxX = location.X
	}
	if location.Y > m.maxY {
		m.maxY = location.Y
	}
	if m.grid == nil {
		m.grid = make(map[Location]TreeHeight)
	}
	m.grid[location] = value
}

func (m *Map) Visible(location Location) bool {
	var x, y uint
	x, y = location.X, location.Y

	if x == 0 || x == m.maxX || y == 0 || y == m.maxY {
		return true // at the edge of the map
	}

	var height TreeHeight
	height = m.Get(location)

	var i uint
	var cursor *Location
	cursor = &Location{}
	for i = 0; i <= x; i++ {
		if i == x {
			return true // visible from the left
		}
		cursor.X = i
		cursor.Y = y
		if m.Get(*cursor) >= height {
			break
		}
	}
	for i = m.maxX; i >= x; i-- {
		if i == x {
			return true // visible from the right
		}
		cursor.X = i
		cursor.Y = y
		if m.Get(*cursor) >= height {
			break
		}
	}
	for i = 0; i <= y; i++ {
		if i == y {
			return true // visible from the top
		}
		cursor.X = x
		cursor.Y = i
		if m.Get(*cursor) >= height {
			break
		}
	}
	for i = m.maxY; i >= y; i-- {
		if i == y {
			return true // visible from the bottom
		}
		cursor.X = x
		cursor.Y = i
		if m.Get(*cursor) >= height {
			break
		}
	}
	return false
}

func ReadMap(filename string) *Map {
	var cursor *Location
	cursor = &Location{0, 0}
	var trees *Map
	trees = &Map{}
	var height int
	var err error
	for line := range ReadLines(filename) {
		cursor.X = 0
		for _, char := range line {
			height, err = strconv.Atoi(string(char))
			if err != nil {
				log.Fatalf("could not parse tree height: %s", string(char))
			}
			trees.Set(*cursor, TreeHeight(height))
			cursor.X++
		}
		cursor.Y++
	}
	return trees
}

func part1(filename string) string {
	trees := ReadMap(filename)
	var result int
	for location := range trees.grid {
		if trees.Visible(location) {
			result++
		}
	}
	return strconv.Itoa(result)
}

func part2(filename string) string {
	return ""
}
