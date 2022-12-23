package main

import (
	"fmt"
)

type Shape struct {
	filled map[Point]bool
}

func (s *Shape) FromFile(filename string) {
	var line string
	var point Point
	for line = range ReadLines(filename) {
		point.Parse(line)
		s.Add(point)
	}
}

func (s *Shape) Add(p Point) {
	if s.filled == nil {
		s.filled = make(map[Point]bool)
	}
	s.filled[p] = true
}

func (s *Shape) SurfaceArea() int {
	var total int
	var p Point
	var direction Direction
	for p = range s.filled {
		for _, direction = range Neighbors {
			if !s.filled[p.Look(direction)] {
				total++
			}
		}
	}
	return total
}

func part1(filename string) string {
	shape := &Shape{}
	shape.FromFile(filename)
	return fmt.Sprintf("%d", shape.SurfaceArea())
}

func part2(filename string) string {
	return ""
}
