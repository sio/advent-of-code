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

func (m *Map) Exists(location Location) bool {
	_, exists := m.grid[location]
	return exists
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

func (m *Map) ScenicScore(location Location) uint {
	var score uint
	score = 1

	var cursor Location
	cursor = location

	var height TreeHeight
	height = m.Get(location)

	var i, seen uint
	seen = 0
	for i = location.X + 1; i <= m.maxX; i++ {
		cursor.X = i
		cursor.Y = location.Y
		if !m.Exists(cursor) {
			break
		}
		seen++
		if m.Get(cursor) >= height {
			break
		}
	}
	score *= seen

	seen = 0
	for i = location.X - 1; i >= 0; i-- {
		cursor.X = i
		cursor.Y = location.Y
		if !m.Exists(cursor) {
			break
		}
		seen++
		if m.Get(cursor) >= height {
			break
		}
	}
	score *= seen

	seen = 0
	for i = location.Y + 1; i <= m.maxY; i++ {
		cursor.X = location.X
		cursor.Y = i
		if !m.Exists(cursor) {
			break
		}
		seen++
		if m.Get(cursor) >= height {
			break
		}
	}
	score *= seen

	seen = 0
	for i = location.Y - 1; i >= 0; i-- {
		cursor.X = location.X
		cursor.Y = i
		if !m.Exists(cursor) {
			break
		}
		seen++
		if m.Get(cursor) >= height {
			break
		}
	}
	score *= seen

	return score
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
	trees := ReadMap(filename)
	var max, current uint
	for location := range trees.grid {
		current = trees.ScenicScore(location)
		if current > max {
			max = current
		}
	}
	return strconv.Itoa(int(max))
}
