package main

import (
	"fmt"
	"math"
	"strings"
)

const (
	cubeFaces     = 6
	cubeFaceEdges = 4
)

type Cube struct {
	face map[Point]*CubeFace
	size Coordinate
}

func (cube *Cube) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Cube(size=%d):", cube.size))
	for _, face := range cube.face {
		b.WriteString("\n  ")
		b.WriteString(face.String())
	}
	return b.String()
}

func (cube *Cube) Face(p Point) *CubeFace {
	p.X /= cube.size
	p.Y /= cube.size
	return cube.face[p]
}

func (cube *Cube) addFace(corner Point) *CubeFace {
	face, ok := cube.face[corner]
	if ok {
		return face
	}
	cube.face[corner] = &CubeFace{
		corner: corner,
	}
	return cube.face[corner]
}

func (cube *Cube) Parse(maze *Maze) {
	cursor := maze.player.location
	cube.face = make(map[Point]*CubeFace)
	cube.size = Coordinate(math.Sqrt(float64(len(maze.tile) / cubeFaces)))

	near := make([]Step, cubeFaceEdges)
	near[Right] = Step{cube.size, 0}
	near[Left] = Step{-cube.size, 0}
	near[Up] = Step{0, -cube.size}
	near[Down] = Step{0, cube.size}

	cube.addFace(cursor)
	for cube.Validate() != nil {
		var change bool = false
		for _, face := range cube.face {
		neighbor_loop:
			for index, neighbor := range face.neighbor {
				if neighbor != nil {
					continue
				}
				side := Facing(index)
				cursor = face.corner.Move(near[side])

				// Simple case: cube faces are laid out near each other in maze
				if maze.Contains(cursor) {
					face.AddNeighbor(side, cube.addFace(cursor), side)
					change = true
					continue neighbor_loop
				}

				// Slightly trickier: adjacent faces are touching via a single corner only
				for _, rotation := range []Rotation{Clockwise, CounterClockwise} {
					facing := side.Turn(rotation)
					adjacent := cursor.Move(near[facing])
					if !maze.Contains(adjacent) {
						continue
					}
					face.AddNeighbor(side, cube.addFace(adjacent), facing)
					change = true
					continue neighbor_loop
				}

				// Not adjacent in maze at all: use what we know about the cube so far
				for _, rotation := range []Rotation{Clockwise, CounterClockwise} {
					detour := side.Turn(rotation)
					interim := face.neighbor[detour]
					if interim == nil {
						continue
					}
					interimOut := face.into[detour].Turn(!rotation)
					neighbor := interim.neighbor[interimOut]
					if neighbor == nil {
						continue
					}
					face.AddNeighbor(side, neighbor, interim.into[interimOut].Turn(rotation))
					change = true
					continue neighbor_loop
				}
			}
		}
		if !change {
			fmt.Println(cube)
			fmt.Println(cube.Validate())
			panic("parser has entered an endless loop")
		}
	}
}

func (cube *Cube) Validate() error {
	if len(cube.face) != cubeFaces {
		return fmt.Errorf("cube has %d faces instead of %d", len(cube.face), cubeFaces)
	}
	for _, face := range cube.face {
		if face == nil {
			return fmt.Errorf("face %v is nil", face)
		}
		seen := make(map[Point]bool)
		for side, neighbor := range face.neighbor {
			if neighbor == nil {
				return fmt.Errorf("%v: neighbor #%d is nil", face, side)
			}
			if neighbor.corner == face.corner || neighbor == face {
				return fmt.Errorf("face references itself as its neighbor: %v", face)
			}
			if seen[neighbor.corner] {
				return fmt.Errorf("same neighbor on multiple sides: %v", face)
			}
			seen[neighbor.corner] = true
		}
	}
	return nil
}

type CubeFace struct {
	corner   Point // top-left corner of cube face
	neighbor [cubeFaceEdges]*CubeFace
	into     [cubeFaceEdges]*Facing // new facing after steping into neighbor
}

func (face *CubeFace) String() string {
	var b strings.Builder
	for i := 0; i < len(face.neighbor); i++ {
		neighbor := face.neighbor[i]
		into := face.into[i]
		b.WriteString(Facing(i).String())
		if neighbor != nil {
			b.WriteString(fmt.Sprint(neighbor.corner))
		} else {
			b.WriteString("{nil}")
		}
		if into != nil {
			b.WriteString(fmt.Sprint(*into))
		} else {
			b.WriteString("?")
		}
		if i+1 < len(face.neighbor) {
			b.WriteString(" ")
		}
	}

	return fmt.Sprintf("Face%v: %s", face.corner, b.String())
}

func (face *CubeFace) AddNeighbor(side Facing, neighbor *CubeFace, into Facing) {
	if face.neighbor[side] != nil {
		return
	}
	face.neighbor[side] = neighbor
	face.into[side] = &into
	fmt.Printf("new forward link\n  from %v\n    to %v\n", face, neighbor)

	backlink := neighbor.neighbor[into.Reverse()]
	switch {
	case backlink == face:
		return
	case backlink != nil:
		panic(fmt.Sprintf("%v: backlink points to a different face already: %v", face, backlink))
	}
	neighbor.neighbor[into.Reverse()] = face
	back := side.Reverse()
	neighbor.into[into.Reverse()] = &back
	fmt.Printf("new backward link\n  from %v\n    to %v\n", neighbor, face)
}
