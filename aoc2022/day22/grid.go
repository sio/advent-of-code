package main

import (
	"fmt"
)

type Facing int

const (
	Right Facing = 0
	Down         = 1
	Left         = 2
	Up           = 3
)

func (f Facing) Turn(r Rotation) Facing {
	if r == Clockwise {
		f++
	} else {
		f--
	}
	if f < Right {
		f = Up
	}
	if f > Up {
		f = Right
	}
	return f
}

func (f Facing) Reverse() Facing {
	return (f + 2) % 4
}

func (f Facing) String() string {
	symbols := [...]string{">", "v", "<", "^"}
	return symbols[f]
}

type Rotation bool

const (
	Clockwise        Rotation = true
	CounterClockwise          = false
)

type Player struct {
	location Point
	facing   Facing
}

func (p *Player) Turn(r Rotation) {
	p.facing = p.facing.Turn(r)
}

func (p *Player) Password() int {
	return int(p.location.Y)*1000 + int(p.location.X)*4 + int(p.facing)
}

func (p *Player) Ahead() Point {
	var next Point
	next = p.location
	switch p.facing {
	default:
		panic(fmt.Sprintf("facing unknown direction: %v", p.facing))
	case Up:
		next.Y--
	case Down:
		next.Y++
	case Left:
		next.X--
	case Right:
		next.X++
	}
	return next
}
