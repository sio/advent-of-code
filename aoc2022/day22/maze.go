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

func (p Point) Next(to Facing) Point {
	switch to {
	default:
		panic(fmt.Sprintf("facing unknown direction: %v", to))
	case Up:
		p.Y--
	case Down:
		p.Y++
	case Left:
		p.X--
	case Right:
		p.X++
	}
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
	cube       *Cube
}

func (m *Maze) Load(filename string) {
	var iter LineIterator
	if err := iter.Open(filename); err != nil {
		panic(err)
	}
	defer func() {
		if err := iter.Close(); err != nil {
			panic(err)
		}
	}()

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

func (m *Maze) ParseCube() {
	m.cube = &Cube{}
	m.cube.Parse(m)
}

func (m *Maze) Step() (ok bool) {
	var current, next Point
	current = m.player.location
	next = m.player.Ahead()

	var facing Facing
	facing = m.player.facing

	if !m.Contains(next) {
		if m.cube != nil {
			// Cube cutout navigation
			next, facing = m.cube.Next(current, m.player.facing)
		} else {
			// Plain maze navigation
			if current.X == next.X { // same column
				if current.Y > next.Y { // going up and wrapping
					next.Y = m.col[next.X].Max
				} else { // going down and wrapping
					next.Y = m.col[next.X].Min
				}
			} else { // same row
				if current.X > next.X { // going left and wrapping
					next.X = m.row[next.Y].Max
				} else { // going right and wrapping
					next.X = m.row[next.Y].Min
				}
			}
		}
	}

	if next == current {
		panic(fmt.Sprintf("move without location change: %v", m.player))
	}
	if !m.Contains(next) {
		panic(fmt.Sprintf("wandering outside of the maze: %v", next))
	}

	if m.tile[next] == Wall {
		return false
	}
	m.player.location = next
	m.player.facing = facing
	return true
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
