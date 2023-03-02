package main

import (
	"fmt"
	"os"
	"strings"
)

var Moves = [...]Direction{
	Down,
	Right,
	Up,
	Left,
}

type BlizzardBasin struct {
	width, height  ScaleUnit
	entrance, exit Point
	blizzard       []Blizzard
	wall           PointSet
}

type Blizzard struct {
	spawn     Point
	direction Direction
}

func (b Blizzard) String() string {
	return fmt.Sprintf("Blizzard{%d,%d,%s}", b.spawn.X, b.spawn.Y, b.direction)
}

func (bb *BlizzardBasin) String() string {
	return fmt.Sprintf(
		"BlizzardBasin(width=%d, height=%d, blizzards=%d, entrance=%v, exit=%v)",
		bb.width,
		bb.height,
		len(bb.blizzard),
		bb.entrance,
		bb.exit,
	)
}

func (bb *BlizzardBasin) Load(filename string) {
	var input []byte
	var err error
	input, err = os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var tile byte
	var cursor Point
	bb.wall = make(PointSet)
	for _, tile = range input {
		if bb.width != 0 && cursor.X > bb.width {
			panic(fmt.Sprintf("unexpected long line at cursor %v, rectangular basin was assumed", cursor))
		}
		direction, isBlizzard := iconDirection[tile]
		switch {
		case tile == byte('\n'):
			cursor.Y++
			if bb.width == 0 {
				bb.width = cursor.X
			}
			cursor.X = 0
			continue
		case tile == byte('#'):
			bb.wall.Add(cursor)
		case tile == byte('.'):
			bb.exit = cursor
			if bb.entrance.X == 0 && bb.entrance.Y == 0 {
				bb.entrance = cursor
			}
		case isBlizzard:
			bb.blizzard = append(bb.blizzard, Blizzard{
				spawn:     cursor,
				direction: direction,
			})
		default:
			panic(fmt.Sprintf("unhandled byte at position %v: %s", cursor, string(tile)))
		}
		cursor.X++
	}
	bb.height = bb.exit.Y - bb.entrance.Y + 1
}

func (bb *BlizzardBasin) Blizzards(round int) PointSet {
	locations := make(PointSet)
	for _, blizzard := range bb.blizzard {
		var slots, offset, sign, dest ScaleUnit
		if blizzard.direction == Up || blizzard.direction == Down {
			slots = bb.height - 2 // wall positions are not available for blizzard placement
		} else {
			slots = bb.width - 2
		}

		switch blizzard.direction {
		case Down:
			offset = blizzard.spawn.Y - 1
			sign = 1
		case Up:
			offset = (bb.height - 2) - (blizzard.spawn.Y - 1)
			sign = -1
		case Right:
			offset = blizzard.spawn.X - 1
			sign = 1
		case Left:
			offset = (bb.width - 2) - (blizzard.spawn.X - 1)
			sign = -1
		default:
			panic("invalid blizzard direction")
		}

		dest = (sign*((ScaleUnit(round)+offset)%slots)+slots)%slots + 1
		var next Point
		switch blizzard.direction {
		case Up, Down:
			next = Point{
				X: blizzard.spawn.X,
				Y: dest,
			}
		case Left, Right:
			next = Point{
				Y: blizzard.spawn.Y,
				X: dest,
			}
		}
		locations.Add(next)
		//fmt.Printf("%v at location %v [%t]\n", blizzard, next, blizzard.spawn == next)
	}
	return locations
}

func (bb *BlizzardBasin) Render(round int) string {
	blizzards := bb.Blizzards(round)

	var builder strings.Builder
	var cursor Point
	for cursor.Y = 0; cursor.Y < bb.height; cursor.Y++ {
		for cursor.X = 0; cursor.X < bb.width; cursor.X++ {
			var isWall, isBlizzard bool
			_, isWall = bb.wall[cursor]
			_, isBlizzard = blizzards[cursor]
			if isWall && isBlizzard {
				panic(fmt.Sprintf("blizzard collided into wall at: %v", cursor))
			}
			if isWall {
				builder.WriteRune('#')
				continue
			}
			if isBlizzard {
				builder.WriteRune('X')
				continue
			}
			builder.WriteRune('.')
		}
		builder.WriteRune('\n')
	}
	return builder.String()
}
