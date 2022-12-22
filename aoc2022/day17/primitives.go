package main

// Point{0, 0} is at the bottom left corner
type Point struct {
	X, Y int64
}

type Direction Point

func (p *Point) Move(d Direction) {
	p.X += d.X
	p.Y += d.Y
}

func (p *Point) Peek(d Direction) Point {
	return Point{
		p.X + d.X,
		p.Y + d.Y,
	}
}

var (
	Down  = Direction{0, -1}
	Left  = Direction{-1, 0}
	Right = Direction{1, 0}
)
