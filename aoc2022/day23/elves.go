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
