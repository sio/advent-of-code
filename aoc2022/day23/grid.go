package main

type Coordinate int

type Direction Point

var (
	North     = Direction{0, -1}
	South     = Direction{0, 1}
	East      = Direction{1, 0}
	West      = Direction{-1, 0}
	NorthWest = Direction{-1, -1}
	SouthWest = Direction{-1, 1}
	NorthEast = Direction{1, -1}
	SouthEast = Direction{1, 1}
)

var Perimeter = [...]Direction{
	North,
	NorthEast,
	East,
	SouthEast,
	South,
	SouthWest,
	West,
	NorthWest,
}

type Point struct {
	X, Y Coordinate
}

func (p Point) Look(d Direction) Point {
	p.X += d.X
	p.Y += d.Y
	return p
}
