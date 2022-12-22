package main

import (
	"fmt"
	"strings"
)

const ChamberWidth = 7

type SkyLine [ChamberWidth]int64

type Chamber struct {
	base           int64
	height         int64
	width          int64
	rocks          map[Point]bool
	spawnFrom      []Shape
	spawnCounter   int
	pushDirections []Direction
	pushCounter    int
	floor          int64
	skyline        SkyLine
	seen           map[ChamberSnapshot]ChamberStatus
	loop           ChamberLoop
}

type ChamberSnapshot struct {
	shape   int
	wind    int
	skyline SkyLine
}

type ChamberStatus struct {
	iteration int64
	height    int64
}

type ChamberLoop struct {
	base  ChamberStatus
	delta ChamberStatus
}

func (chamber *Chamber) Height() int64 {
	return chamber.base + chamber.height
}

func (chamber *Chamber) DropN(iterations int64) {
	var i int64
	for i = 0; i < iterations; i++ {
		if !chamber.Next(i) {
			break
		}
	}

	var remaining int64
	remaining = iterations - 1 - i
	if remaining <= 0 {
		return // did not encounter a loop
	}

	chamber.base = (remaining / chamber.loop.delta.iteration) * chamber.loop.delta.height
	for i = 0; i < remaining%chamber.loop.delta.iteration; i++ {
		chamber.Next(-1)
	}
}

func (chamber *Chamber) Next(iteration int64) bool {
	var shape *Shape
	shape = chamber.Spawn()
	for {
		chamber.Push(shape)
		if !chamber.Descend(shape) {
			break
		}
	}
	chamber.Settle(shape)

	if iteration < 0 {
		return false // do not update skyline and floor and the rest of parameters
	}

	chamber.UpdateSkyLine()
	//chamber.Cut() // detecting loops make us so fast we don't need to cut old data

	var state = ChamberSnapshot{
		shape:   chamber.spawnCounter,
		wind:    chamber.pushCounter,
		skyline: chamber.skyline,
	}
	if chamber.seen == nil {
		chamber.seen = make(map[ChamberSnapshot]ChamberStatus)
	}
	var found bool
	var prev ChamberStatus
	prev, found = chamber.seen[state]
	if !found {
		chamber.seen[state] = ChamberStatus{
			height:    chamber.height,
			iteration: iteration,
		}
		return true // continue iteration
	}
	chamber.loop = ChamberLoop{
		base: ChamberStatus{
			iteration: iteration,
			height:    chamber.height,
		},
		delta: ChamberStatus{
			iteration: iteration - prev.iteration,
			height:    chamber.height - prev.height,
		},
	}
	return false // found a loop, no need to continue
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

func (chamber *Chamber) UpdateSkyLine() {
	var x, y int64
	for x = 0; x < chamber.width; x++ {
		for y = chamber.height; y > 0; y-- {
			if chamber.rocks[Point{x, y}] {
				chamber.skyline[x] = chamber.height - y
				break
			}
		}
	}
}

func (chamber *Chamber) Cut() {
	var floor int64
	floor = chamber.height

	var x int64
	for x = 0; x < chamber.width; x++ {
		if chamber.height-chamber.skyline[x] < floor {
			floor = chamber.height - chamber.skyline[x]
		}
	}

	var rock Point
	for rock = range chamber.rocks {
		if rock.Y < floor {
			delete(chamber.rocks, rock)
		}
	}
	chamber.floor = floor
}

func (chamber *Chamber) Render() string {
	var builder strings.Builder
	var x, y int64
	for y = chamber.height - 1; y >= chamber.floor; y-- {
		for x = 0; x < chamber.width; x++ {
			if chamber.rocks[Point{x, y}] {
				builder.WriteRune(Rock)
			} else {
				builder.WriteRune(Air)
			}
		}
		builder.WriteRune('\n')
	}
	if chamber.floor > 0 {
		builder.WriteString(fmt.Sprintf("floor at %d\n", chamber.floor))
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

func Play(filename string, rounds int64) string {
	chamber := Chamber{width: ChamberWidth}
	chamber.ReadJetPattern(filename)
	chamber.DropN(rounds)
	return fmt.Sprintf("%d", chamber.Height())
}

func part1(filename string) string {
	return Play(filename, 2022)
}

func part2(filename string) string {
	return Play(filename, 1000000000000)
}
