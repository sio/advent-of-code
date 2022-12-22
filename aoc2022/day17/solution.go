package main

import (
	"fmt"
	"strings"
)

type Chamber struct {
	height int
	width  int
	rocks  map[Point]bool
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
			panic(fmt.Sprintf("collision at {%d,%d}", dest.X, dest.Y))
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
	var x, y int
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

func part1(filename string) string {
	shapes := DefaultShapes()
	chamber := Chamber{width: 7}
	y := 0
	for _, s := range shapes {
		s.corner.Y = y
		chamber.Settle(&s)
		y += s.height + 1
	}
	fmt.Println(chamber.Render())
	return ""
}

func part2(filename string) string {
	return ""
}
