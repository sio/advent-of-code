package main

import (
	"fmt"
)

type Shape struct {
	filled map[Point]bool
	top    Point
}

func (s *Shape) FromFile(filename string) {
	var line string
	var point Point
	for line = range ReadLines(filename) {
		point.Parse(line)
		s.Add(point)
		if point.Z > s.top.Z {
			s.top = point
		}
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

func (s *Shape) ProperArea() int {
	var up = Direction{0, 0, 1}
	return s.SurfaceWalk(s.top.Look(up), nil)
}

func (s *Shape) SurfaceWalk(cursor Point, visited map[Point]bool) int {
	if s.filled[cursor] {
		panic("can walk only outside of the shape")
	}
	if visited == nil {
		visited = make(map[Point]bool)
	}
	visited[cursor] = true
	var direction Direction
	var surface int
	var neighbor Point
	for _, direction = range Neighbors {
		neighbor = cursor.Look(direction)
		if s.filled[neighbor] {
			surface++
			continue
		}
		if visited[neighbor] || !s.Adjacent(neighbor) {
			continue
		}
		surface += s.SurfaceWalk(neighbor, visited)
	}
	return surface
}

// Steam may flow just outside of immediate contact for a moment
func (s *Shape) Adjacent(p Point) bool {
	var direction Direction
	for _, direction = range Neighbors {
		if s.filled[p.Look(direction)] {
			return true
		}
	}
	for _, direction = range Diagonals {
		if s.filled[p.Look(direction)] {
			return true
		}
	}
	return false
}

func part1(filename string) string {
	shape := &Shape{}
	shape.FromFile(filename)
	return fmt.Sprintf("%d", shape.SurfaceArea())
}

func part2(filename string) string {
	shape := &Shape{}
	shape.FromFile(filename)
	return fmt.Sprintf("%d", shape.ProperArea())
}
