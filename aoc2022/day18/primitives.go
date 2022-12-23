package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Point struct {
	X, Y, Z int
}

func (p *Point) Parse(line string) {
	var chunk []string
	var coord [3]int
	chunk = strings.Split(line, ",")
	if len(chunk) != len(coord) {
		panic(fmt.Sprintf("invalid data point: %s", line))
	}
	var i int
	var err error
	for i = 0; i < len(coord); i++ {
		coord[i], err = strconv.Atoi(chunk[i])
		if err != nil {
			panic(fmt.Sprintf("invalid coordinate %s: %s", chunk[i], line))
		}
	}
	p.X = coord[0]
	p.Y = coord[1]
	p.Z = coord[2]
}

type Direction Point

func (p *Point) Look(d Direction) Point {
	return Point{p.X + d.X, p.Y + d.Y, p.Z + d.Z}
}

var Neighbors = []Direction{
	{1, 0, 0},
	{-1, 0, 0},
	{0, 1, 0},
	{0, -1, 0},
	{0, 0, 1},
	{0, 0, -1},
}
