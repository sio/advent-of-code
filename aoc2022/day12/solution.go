package main

import (
	"fmt"
	"log"
	"strconv"
)

type Point struct {
	X, Y int
}

type Step Point

func (p *Point) Move(delta Step) {
	p.X += delta.X
	p.Y += delta.Y
}

func (p *Point) Neighbor(delta Step) Point {
	return Point{
		X: p.X + delta.X,
		Y: p.Y + delta.Y,
	}
}

type Map struct {
	Height   map[Point]rune
	Start    Point
	Finish   Point
	distance map[Point]int
}

func (m *Map) Print() {
	var x, y int
	var height rune
	var exists, empty bool
	var cursor Point
	for x = 0; true; x++ {
		empty = true
		for y = 0; true; y++ {
			cursor = Point{x, y}
			height, exists = m.Height[cursor]
			if !exists {
				break
			}
			switch cursor {
			case m.Start:
				fmt.Print("S")
			case m.Finish:
				fmt.Print("E")
			default:
				fmt.Printf("%c", height)
			}
			empty = false
		}
		if empty {
			break
		}
		fmt.Print("\n")
	}
}

func (m *Map) Route() {
	unvisited := make(map[Point]bool)

	var cursor, pos, next Point
	for pos = range m.Height {
		unvisited[pos] = true
	}

	var dist int
	var found, first bool
	cursor = m.Start
	for {
		m.visit(cursor)
		delete(unvisited, cursor)
		if len(unvisited) == 0 {
			break
		}
		//log.Printf("Positions left to evaluate: %d/%d", len(unvisited), len(m.Height))
		first = true
		for pos = range unvisited { // select the lowest distance among unvisited
			dist, found = m.distance[pos]
			if !found {
				continue
			}
			if first {
				next = pos
				first = false
			}
			if dist < m.distance[next] {
				next = pos
			}
		}
		if next == m.Finish {
			break
		}
		if next == cursor {
			panic(fmt.Sprintf("entering endless loop current=next=(%d,%d)", next.X, next.Y))
		}
		cursor = next
	}
}

func (m *Map) visit(cursor Point) {
	var steps = []Step{
		{X: 0, Y: 1},
		{X: 0, Y: -1},
		{X: 1, Y: 0},
		{X: -1, Y: 0},
	}
	var height rune
	var exists bool
	var step Step
	var neighbor Point
	var newDistance, oldDistance int

	//log.Printf("Visiting point (%d,%d), height %c, distance %d", cursor.X, cursor.Y, m.Height[cursor], m.distance[cursor])
	for _, step = range steps {
		neighbor = cursor.Neighbor(step)
		//fmt.Printf("Neighbor (%d,%d)", neighbor.X, neighbor.Y)
		height, exists = m.Height[neighbor]
		if !exists {
			//fmt.Printf(" does not exist (height=%c)\n", height)
			continue
		}
		//fmt.Printf(" height=%c", height)
		if height-m.Height[cursor] > 1 {
			//fmt.Printf(" no path (too high up)\n")
			continue
		}
		newDistance = m.distance[cursor] + 1
		oldDistance, exists = m.distance[neighbor]
		if !exists || newDistance < oldDistance {
			m.distance[neighbor] = newDistance
		}
		//fmt.Printf(" distance=%d\n", m.distance[neighbor])
	}
}

func NewMap() *Map {
	m := &Map{}
	m.Height = make(map[Point]rune)
	m.distance = make(map[Point]int)
	return m
}

func ParseArea(filename string) (*Map, error) {
	var x, y int
	var line string
	var char rune
	var area = NewMap()
	for line = range ReadLines(filename) {
		x = 0
		for _, char = range line {
			switch char {
			case 'S':
				area.Start = Point{x, y}
				char = 'a'
			case 'E':
				area.Finish = Point{x, y}
				char = 'z'
			}
			if char > 'z' || char < 'a' {
				return nil, fmt.Errorf("invalid area height at position (%d,%d): %q", x, y, char)
			}
			area.Height[Point{x, y}] = char
			x++
		}
		y++
	}
	return area, nil
}

func part1(filename string) string {
	area, err := ParseArea(filename)
	if err != nil {
		log.Fatal(err)
	}
	area.Route()
	return strconv.Itoa(area.distance[area.Finish])
}

func part2(filename string) string {
	return ""
}
