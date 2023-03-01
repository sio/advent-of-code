package main

import (
	"fmt"
	"strings"
)

type void struct{}

const (
	axisMaxValue = Coordinate(^uint(0) >> 1)
	axisMinValue = -axisMaxValue - 1
)

type PointSet map[Point]void

func (set *PointSet) Add(p Point) error {
	if set.Contains(p) {
		return fmt.Errorf("point already in set: %v", p)
	}
	(*set)[p] = void{}
	return nil
}

func (set *PointSet) Contains(p Point) bool {
	_, exist := (*set)[p]
	return exist
}

type ElfGroup struct {
	elves PointSet
	max   Point
	min   Point
}

func (group *ElfGroup) String() string {
	group.updateRectangle()

	var b strings.Builder
	var cursor Point
	for cursor.Y = group.min.Y; cursor.Y <= group.max.Y; cursor.Y++ {
		for cursor.X = group.min.X; cursor.X <= group.max.X; cursor.X++ {
			if group.elves.Contains(cursor) {
				b.WriteRune('#')
			} else {
				b.WriteRune('.')
			}
		}
		b.WriteRune('\n')
	}
	return b.String()
}

func (group *ElfGroup) Load(filename string) {
	var iter LineIterator
	if err := iter.Open(filename); err != nil {
		panic(err)
	}
	defer func() {
		if err := iter.Close(); err != nil {
			panic(err)
		}
	}()

	group.elves = make(PointSet)
	var cursor Point
	for iter.Next() {
		cursor.Y++
		cursor.X = 0
		for _, char := range iter.Value() {
			cursor.X++
			switch char {
			case '.': // noop
			case '#':
				err := group.elves.Add(cursor)
				if err != nil {
					panic(err)
				}
			default:
				panic(fmt.Sprintf("unsupported character at %v: %c", cursor, char))
			}
		}
	}
}

func (group *ElfGroup) updateRectangle() {
	group.min = Point{axisMaxValue, axisMaxValue}
	group.max = Point{axisMinValue, axisMinValue}
	for cursor := range group.elves {
		if cursor.X < group.min.X {
			group.min.X = cursor.X
		}
		if cursor.X > group.max.X {
			group.max.X = cursor.X
		}
		if cursor.Y < group.min.Y {
			group.min.Y = cursor.Y
		}
		if cursor.Y > group.max.Y {
			group.max.Y = cursor.Y
		}
	}
}

var Movements = [...]Direction{
	North,
	South,
	West,
	East,
}

var Row = map[Direction][3]Direction{
	North: {NorthWest, North, NorthEast},
	South: {SouthWest, South, SouthEast},
	West:  {NorthWest, West, SouthWest},
	East:  {NorthEast, East, SouthEast},
}

func (group *ElfGroup) Result() int {
	group.updateRectangle()

	var width, height int
	width = int(group.max.X - group.min.X + 1)
	height = int(group.max.Y - group.min.Y + 1)

	return width*height - len(group.elves)
}

func (group *ElfGroup) Play(rounds int) int {
	var i int
	for i = 0; i < rounds; i++ {
		changed := group.Round(i)
		if !changed {
			break
		}
	}
	return i
}

func (group *ElfGroup) Round(index int) bool {
	moves := make(map[Point]Point) // to -> from
	banned := make(map[Point]bool)

	// Plan movements for each Elf
	for elf := range group.elves {
		var hasNeighbors bool
		for _, direction := range Perimeter {
			if group.elves.Contains(elf.Look(direction)) {
				hasNeighbors = true
				break
			}
		}
		if !hasNeighbors {
			continue
		}
	side_loop:
		for side := 0; side < len(Movements); side++ {
			direction := Movements[(side+index)%len(Movements)]
			for _, neighbor := range Row[direction] {
				if group.elves.Contains(elf.Look(neighbor)) {
					continue side_loop
				}
			}

			destination := elf.Look(direction)

			if banned[destination] {
				break
			}
			_, taken := moves[destination]
			if taken {
				banned[destination] = true
				delete(moves, destination)
				break
			}
			moves[destination] = elf
			break
		}
	}

	var elfCount = len(group.elves)

	// Remove old Elf locations
	for _, from := range moves {
		delete(group.elves, from)
	}

	// Add new Elf locations
	for to, from := range moves {
		if group.elves.Contains(to) {
			err := group.elves.Add(from) // restore Elf at old location
			if err != nil {
				panic(err)
			}
			continue
		}
		err := group.elves.Add(to)
		if err != nil {
			panic(err)
		}
	}

	// Sanity check
	if len(group.elves) != elfCount {
		fmt.Println(group)
		panic(fmt.Sprintf("Elf count changed after movements: was %d, now %d", elfCount, len(group.elves)))
	}

	return len(moves) != 0
}
