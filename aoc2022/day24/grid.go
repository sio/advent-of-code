package main

import (
	"fmt"
)

type ScaleUnit int

type Direction Point

var (
	Up    = Direction{0, -1}
	Down  = Direction{0, 1}
	Right = Direction{1, 0}
	Left  = Direction{-1, 0}
)

var directionIcon = map[Direction]string{
	Up:    "^",
	Down:  "v",
	Left:  "<",
	Right: ">",
}

var iconDirection = map[byte]Direction{
	byte('^'): Up,
	byte('v'): Down,
	byte('>'): Right,
	byte('<'): Left,
}

func (d Direction) String() string {
	icon, found := directionIcon[d]
	if !found {
		return fmt.Sprintf("Direction{%d,%d}", d.X, d.Y)
	}
	return icon
}

type Point struct {
	X, Y ScaleUnit
}

func (p Point) Look(d Direction) Point {
	p.X += d.X
	p.Y += d.Y
	return p
}

type void struct{}

type PointSet map[Point]void

func (set *PointSet) Add(p Point) {
	(*set)[p] = void{}
}

func (set *PointSet) Contains(p Point) bool {
	_, found := (*set)[p]
	return found
}
