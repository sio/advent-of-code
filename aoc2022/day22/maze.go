package main

import (
	"fmt"
	"strings"
)

type Coordinate int

type Boundary struct {
	Min Coordinate
	Max Coordinate
}

type Point struct {
	X, Y Coordinate
}

type Cell uint8

const (
	Undefined Cell = iota
	Empty
	Wall
)

type Maze struct {
	tile       map[Point]Cell
	row        map[Coordinate]Boundary
	col        map[Coordinate]Boundary
	directions string
}

func (m *Maze) Load(filename string) {
	var iter LineIterator
	if err := iter.Open(filename); err != nil {
		panic(err)
	}
	defer iter.Close()

	m.tile = make(map[Point]Cell)
	m.row = make(map[Coordinate]Boundary)
	m.col = make(map[Coordinate]Boundary)
	m.directions = ""

	var cursor Point
	var endOfMap bool
	for iter.Next() {
		cursor.Y++
		cursor.X = 0
		if len(iter.Value()) == 0 {
			endOfMap = true
			continue
		}
		if endOfMap {
			if len(m.directions) != 0 {
				panic("attempting to overwrite directions")
			}
			m.directions = iter.Value()
			continue
		}
		for _, char := range iter.Value() {
			cursor.X++
			switch char {
			default:
				panic(fmt.Sprintf("unsupported map tile: %c", char))
			case '.':
				m.tile[cursor] = Empty
			case '#':
				m.tile[cursor] = Wall
			case ' ': // do nothing for empty tiles
			}
			if m.tile[cursor] == Undefined {
				continue
			}

			rowBounds := m.row[cursor.Y]
			if cursor.X < rowBounds.Min || rowBounds.Min == 0 {
				rowBounds.Min = cursor.X
			}
			if cursor.X > rowBounds.Max {
				rowBounds.Max = cursor.X
			}
			if rowBounds != m.row[cursor.Y] {
				m.row[cursor.Y] = rowBounds
			}

			colBounds := m.col[cursor.X]
			if cursor.Y < colBounds.Min || colBounds.Min == 0 {
				colBounds.Min = cursor.Y
			}
			if cursor.Y > colBounds.Max {
				colBounds.Max = cursor.Y
			}
			if colBounds != m.col[cursor.X] {
				m.col[cursor.X] = colBounds
			}
		}
	}
}

func (m *Maze) String() string {
	var b strings.Builder
	var cursor Point
	for {
		cursor.Y++
		row := m.row[cursor.Y]
		for cursor.X = 1; cursor.X <= row.Max; cursor.X++ {
			var tile rune
			switch m.tile[cursor] {
			case Undefined:
				tile = ' '
			case Empty:
				tile = '.'
			case Wall:
				tile = '#'
			default:
				panic(fmt.Sprintf("unsupported tile value: %d", m.tile[cursor]))
			}
			b.WriteRune(tile)
		}
		if cursor.X < 2 {
			break
		}
		b.WriteRune('\n')
	}
	b.WriteRune('\n')
	b.WriteString(m.directions)
	return b.String()
}
