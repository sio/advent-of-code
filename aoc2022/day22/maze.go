package main

import (
	"fmt"
	"strconv"
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

type Step Point

func (p Point) Move(s Step) Point {
	p.X += s.X
	p.Y += s.Y
	return p
}

func (a Point) Distance(b Point) Step {
	return Step{
		X: b.X - a.X,
		Y: b.Y - a.Y,
	}
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
	player     Player
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

	m.player = Player{
		location: Point{m.row[1].Min, 1},
		facing:   Right,
	}
}

func (m *Maze) Step() (ok bool) {
	var from, to Point
	from = m.player.location
	to = m.player.Ahead()

	var tile Cell
	var found bool
	tile, found = m.tile[to]
	if !found { // wrap around the map
		if from.X == to.X { // same column
			if from.Y > to.Y { // going up and wrapping
				to.Y = m.col[to.X].Max
			} else { // going down and wrapping
				to.Y = m.col[to.X].Min
			}
		} else { // same row
			if from.X > to.X { // going left and wrapping
				to.X = m.row[to.Y].Max
			} else { // going right and wrapping
				to.X = m.row[to.Y].Min
			}
		}
		tile = m.tile[to]
	}

	switch tile {
	case Wall:
		return false // can not move forward
	case Empty:
		m.player.location = to
		return true
	default:
		panic(fmt.Sprintf("unsupported tile: %v", tile))
	}
}

func (m *Maze) Play() {
	var buf strings.Builder
	for _, char := range m.directions {
		var rotation Rotation
		switch char {
		case 'R':
			rotation = Clockwise
		case 'L':
			rotation = CounterClockwise
		default:
			buf.WriteRune(char)
			continue
		}
		m.run(buf.String())
		buf.Reset()
		m.player.Turn(rotation)
	}
	m.run(buf.String())
}

func (m *Maze) run(far string) {
	steps, err := strconv.Atoi(far)
	if err != nil {
		panic(err)
	}
	for i := 0; i < steps; i++ {
		if !m.Step() {
			break
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
			if cursor == m.player.location {
				switch m.player.facing {
				case Up:
					tile = '^'
				case Down:
					tile = 'v'
				case Left:
					tile = '<'
				case Right:
					tile = '>'
				default:
					panic(fmt.Sprintf("player is facing in invalid direction: %v", m.player.facing))
				}
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

func (m *Maze) Contains(p Point) bool {
	return m.tile[p] != Undefined
}
