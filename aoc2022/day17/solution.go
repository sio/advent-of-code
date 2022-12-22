package main

import (
	"fmt"
	"strings"
)

type Chamber struct {
	height         int64
	width          int64
	rocks          map[Point]bool
	spawnFrom      []Shape
	spawnCounter   int
	pushDirections []Direction
	pushCounter    int
}

func (chamber *Chamber) Drop() {
	var shape *Shape
	shape = chamber.Spawn()
	for {
		chamber.Push(shape)
		if !chamber.Descend(shape) {
			break
		}
	}
	chamber.Settle(shape)
}

func (chamber *Chamber) Descend(shape *Shape) bool {
	var rock Direction
	var corner Point
	corner = shape.corner.Peek(Down)
	if corner.Y < 0 {
		return false
	}
	for _, rock = range shape.rocks {
		if chamber.rocks[corner.Peek(rock)] {
			return false
		}
	}
	shape.corner.Move(Down)
	return true
}

func (chamber *Chamber) Push(shape *Shape) {
	if len(chamber.pushDirections) == 0 {
		panic("chamber push directions not initialized")
	}

	var direction Direction
	direction = chamber.pushDirections[chamber.pushCounter]
	chamber.pushCounter++
	chamber.pushCounter = chamber.pushCounter % len(chamber.pushDirections)

	if direction.Y != 0 {
		panic(fmt.Sprintf("pushing is allowed only along horizontal axis: %v", direction))
	}

	var corner Point
	corner = shape.corner.Peek(direction)
	if corner.X < 0 || corner.X+shape.width > chamber.width {
		return // don't go outside of chamber walls
	}
	var rock Direction
	for _, rock = range shape.rocks {
		if chamber.rocks[corner.Peek(rock)] {
			return // don't allow collision with previously settled rocks
		}
	}
	shape.corner.Move(direction)
}

func (chamber *Chamber) Spawn() *Shape {
	if len(chamber.spawnFrom) == 0 {
		chamber.spawnFrom = DefaultShapes()
	}
	var shape Shape
	shape = chamber.spawnFrom[chamber.spawnCounter]
	shape.corner.X = 2
	shape.corner.Y = 3 + chamber.height
	chamber.spawnCounter++
	chamber.spawnCounter = chamber.spawnCounter % len(chamber.spawnFrom)
	return &shape
}

func (chamber *Chamber) Settle(shape *Shape) {
	if chamber.rocks == nil {
		chamber.rocks = make(map[Point]bool)
	}

	var rock Direction
	var dest Point
	for _, rock = range shape.rocks {
		dest = shape.corner.Peek(rock)
		if chamber.rocks[dest] {
			panic(fmt.Sprintf("collision at %v: shape %v", dest, *shape))
		}
		if 0 > dest.X || dest.X >= chamber.width {
			panic(fmt.Sprintf("placing a rock outside of chamber at {%d,%d}", dest.X, dest.Y))
		}
		chamber.rocks[dest] = true
		if dest.Y+1 > chamber.height {
			chamber.height = dest.Y + 1
		}
	}
}

func (chamber *Chamber) Render() string {
	var builder strings.Builder
	var x, y int64
	for y = chamber.height - 1; y >= 0; y-- {
		for x = 0; x < chamber.width; x++ {
			if chamber.rocks[Point{x, y}] {
				builder.WriteRune(Rock)
			} else {
				builder.WriteRune(Air)
			}
		}
		builder.WriteRune('\n')
	}
	return builder.String()
}

func (chamber *Chamber) ReadJetPattern(filename string) {
	var char rune
	var direction Direction
	for char = range ReadChars(filename) {
		switch char {
		case '<':
			direction = Left
		case '>':
			direction = Right
		case '\n':
			continue
		case '\r':
			continue
		default:
			panic(fmt.Sprintf("unsupported direction: %c (%v)", char, char))
		}
		chamber.pushDirections = append(chamber.pushDirections, direction)
	}
}

func part1(filename string) string {
	chamber := Chamber{width: 7}
	chamber.ReadJetPattern(filename)
	for i := 0; i < 2022; i++ {
		chamber.Drop()
	}
	return fmt.Sprintf("%d", chamber.height)
}

func part2(filename string) string {
	return ""
}
