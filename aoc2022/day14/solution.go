package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Point struct {
	X, Y int
}

func (p *Point) Move(d Direction) {
	p.X += d.X
	p.Y += d.Y
}

func (p *Point) Neighbor(d Direction) Point {
	return Point{p.X + d.X, p.Y + d.Y}
}

func (p *Point) Parse(line string) (err error) {
	coord := strings.Split(line, ",")
	if len(coord) != 2 {
		return fmt.Errorf("invalid coordinates: %s", line)
	}
	p.X, err = strconv.Atoi(coord[0])
	if err != nil {
		return fmt.Errorf("invalid X coordinate (%s): %w", line, err)
	}
	p.Y, err = strconv.Atoi(coord[1])
	if err != nil {
		return fmt.Errorf("invalid Y coordinate (%s): %w", line, err)
	}
	return nil
}

type Direction Point

type Rectangle struct {
	TopLeft     Point
	BottomRight Point
}

func (r *Rectangle) Contains(p Point) bool {
	return r.TopLeft.X <= p.X && r.TopLeft.Y <= p.Y && r.BottomRight.X >= p.X && r.BottomRight.Y >= p.Y
}

type Tile uint8

const Air Tile = 0
const (
	Rock Tile = 1 << iota
	Sand
)

type Map struct {
	tiles  map[Point]Tile
	area   *Rectangle
	recent Point
	floor  int
}

func (m *Map) Draw() string {
	return m.DrawRectangle(m.area)
}

func (m *Map) DrawRectangle(r *Rectangle) string {
	symbols := map[Tile]rune{
		Air:  '.',
		Rock: '#',
		Sand: 'o',
	}
	var build strings.Builder
	var x, y int
	for y = r.TopLeft.Y; y <= r.BottomRight.Y; y++ {
		for x = r.TopLeft.X; x <= r.BottomRight.X; x++ {
			build.WriteRune(symbols[m.Read(Point{x, y})])
		}
		build.WriteString("\n")
	}
	return build.String()
}

func (m *Map) DrawRecent(size int) string {
	recent := &Rectangle{
		TopLeft:     Point{m.recent.X - size/2, m.recent.Y - size/2},
		BottomRight: Point{m.recent.X + size/2, m.recent.Y + size/2},
	}
	return m.DrawRectangle(recent)
}

func (m *Map) Read(place Point) Tile {
	if place.Y == m.floor {
		return Rock
	}
	return m.tiles[place]
}

func (m *Map) Fill(place Point, value Tile) {
	if m.area == nil {
		m.area = &Rectangle{place, place}
	}
	if place.X < m.area.TopLeft.X {
		m.area.TopLeft.X = place.X
	}
	if place.Y < m.area.TopLeft.Y {
		m.area.TopLeft.Y = place.Y
	}
	if place.X > m.area.BottomRight.X {
		m.area.BottomRight.X = place.X
	}
	if place.Y > m.area.BottomRight.Y {
		m.area.BottomRight.Y = place.Y
	}
	m.tiles[place] = value
	m.recent = place
}

func (m *Map) AddFloor(delta int) {
	if m.floor != 0 {
		return
	}
	m.floor = m.area.BottomRight.Y + 2
	m.area.BottomRight.Y = m.floor
}

func (m *Map) PourSand(from Point) (count int) {
	var ok bool
	for {
		ok = m.DropSand(from)
		if !ok {
			break
		}
		count++
		if false { // debug
			fmt.Printf("\n--- Sandbag #%d ---\n%s\n", count, m.Draw())
		}
	}
	return count
}

func (m *Map) DropSand(from Point) (ok bool) {
	m.Fill(from, m.tiles[from]) // extend map area to include sand source
	if m.tiles[from] != Air {
		return false
	}

	steps := []Direction{
		{0, 1},
		{-1, 1},
		{1, 1},
	}

	var sand Point
	sand = from

	var moved bool
	for {
		for _, step := range steps {
			if m.Read(sand.Neighbor(step)) == Air {
				sand.Move(step)
				moved = true
				if m.floor == 0 && !m.area.Contains(sand) {
					return false
				}
				break
			}
			moved = false
		}
		if !moved {
			break
		}
	}
	m.Fill(sand, Sand)
	return true
}

func (m *Map) FromFile(filename string) (err error) {
	if m.tiles == nil {
		m.tiles = make(map[Point]Tile)
	}
	var line string
	var points []string
	var i int
	var start, end, cursor Point
	var direction Direction
	for line = range ReadLines(filename) {
		points = strings.Split(line, " -> ")
		for i = 0; i < len(points)-1; i++ {
			err = start.Parse(points[i])
			if err != nil {
				return err
			}
			err = end.Parse(points[i+1])
			if err != nil {
				return err
			}
			direction = Direction{
				X: sign(end.X - start.X),
				Y: sign(end.Y - start.Y),
			}
			if direction.X != 0 && direction.Y != 0 {
				return fmt.Errorf("rocks must go in horizontal/vertical lines only: %s -> %s", points[i], points[i+1])
			}
			cursor = start
			for {
				m.Fill(cursor, Rock)
				if cursor == end {
					break
				}
				cursor.Move(direction)
			}
		}
	}
	return nil
}

func sign(value int) int {
	switch {
	case value < 0:
		return -1
	case value == 0:
		return 0
	case value > 0:
		return 1
	default:
		panic("impossible branching")
	}
}

func ReadCave(filename string) (*Map, error) {
	cave := &Map{}

	var err error
	err = cave.FromFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not parse file: %v\n", err)
	}
	return cave, nil
}

func part1(filename string) string {
	cave, err := ReadCave(filename)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return strconv.Itoa(cave.PourSand(Point{500, 0}))
}

func part2(filename string) string {
	cave, err := ReadCave(filename)
	if err != nil {
		log.Fatalf("%v", err)
	}
	cave.AddFloor(2)
	return strconv.Itoa(cave.PourSand(Point{500, 0}))
}
