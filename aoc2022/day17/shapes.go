package main

import (
	"strings"
)

const (
	Rock = '#'
	Air  = '.'
)

type Shape struct {
	corner Point // bottom left corner of a rectangle confining the shape
	height int
	width  int
	rocks  []Direction
}

func (s *Shape) Parse(visual string) {
	var lines []string
	lines = strings.Split(visual, "\n")

	var line string
	var char rune
	var x, y int
	for y = 0; y < len(lines); y++ {
		line = lines[len(lines)-(y+1)]
		line = strings.TrimSpace(line)
		for x, char = range line {
			if char != Rock {
				continue
			}
			s.rocks = append(s.rocks, Direction{x, y})
			if x+1 > s.width {
				s.width = x + 1
			}
			if y+1 > s.height {
				s.height = y + 1
			}
		}
	}
}

func DefaultShapes() (shapes []Shape) {
	const raw = `
		####

		.#.
		###
		.#.

		..#
		..#
		###

		#
		#
		#
		#

		##
		## `
	var input []string
	input = strings.Split(raw, "\n\n")

	var i int
	var s Shape
	for i = 0; i < len(input); i++ {
		s = Shape{}
		s.Parse(input[i])
		shapes = append(shapes, s)
	}
	return shapes
}
