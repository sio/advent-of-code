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

func (p *Point) Ray(d Direction) Arrow {
	return Arrow{
		Base: Point{
			X: base(p.X, d.X),
			Y: base(p.Y, d.Y),
			Z: base(p.Z, d.Z),
		},
		Aim: d,
	}
}

func base(coordinate, direction int) int {
	switch {
	case direction == Stay:
		return coordinate
	case direction == Increase:
		return 0
	case direction == Decrease:
		return 0
	default:
		panic("non-adjacent directions support is not implemented")
	}
}

const (
	Increase = 1
	Decrease = -1
	Stay     = 0
)

var Neighbors = []Direction{
	{Increase, 0, 0},
	{Decrease, 0, 0},
	{0, Increase, 0},
	{0, Decrease, 0},
	{0, 0, Increase},
	{0, 0, Decrease},
}

// Unbounded directed line
type Arrow struct {
	// Position where Arrow crosses one of the axis
	Base Point

	// Direction from any Arrow point to the next one
	Aim Direction

	debug bool
}

type RelativePosition int8

const (
	Before RelativePosition = -1
	Equal                   = 0
	After                   = 1
)

// Check whether point A comes before or after point B on this ray
func (ray *Arrow) Compare(a, b Point) RelativePosition {
	var coord = [...]coordComparer{
		{ray.Base.X, ray.Aim.X, a.X, b.X},
		{ray.Base.Y, ray.Aim.Y, a.Y, b.Y},
		{ray.Base.Z, ray.Aim.Z, a.Z, b.Z},
	}
	var cmp, result RelativePosition
	var same int

	for _, c := range coord {
		cmp = c.Compare()
		if ray.debug {
			fmt.Printf("%v: %v\n", c, cmp)
		}
		if cmp == Equal {
			same++
		} else {
			result = cmp
		}
	}
	if ray.debug {
		fmt.Println("")
	}
	if same < 2 {
		panic(fmt.Sprintf("at least two coordinates must be equal: ray=%v, a=%v, b=%v", *ray, a, b))
	}
	return result
}

// Comparison data for one axis
type coordComparer struct {
	rayBase int
	rayAim  int
	coordA  int
	coordB  int
}

func (cc *coordComparer) Compare() RelativePosition {
	//fmt.Printf("Compare() on %v: ", *cc)
	if cc.rayAim == Stay && (cc.rayBase != cc.coordA || cc.rayBase != cc.coordB) {
		panic("one or both points do not belong to this ray")
	}
	if cc.rayAim != Stay && cc.rayAim != Increase && cc.rayAim != Decrease {
		panic("aiming not implemented for non-trivial cases")
	}
	if cc.coordA < cc.coordB {
		return Before * RelativePosition(cc.rayAim)
	}
	if cc.coordA > cc.coordB {
		return After * RelativePosition(cc.rayAim)
	}
	if cc.coordA == cc.coordB {
		return Equal * RelativePosition(cc.rayAim)
	}
	panic("impossible branching")
}
