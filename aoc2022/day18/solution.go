package main

import (
	"fmt"
)

type SolidShape struct {
	boundary map[Arrow]Point
}

func (s *SolidShape) FromFile(filename string) {
	var line string
	var point Point
	for line = range ReadLines(filename) {
		point.Parse(line)
		s.Attach(point)
	}
}

func (s *SolidShape) Verify(filename string) {
	var line string
	var point Point

	correct := make(map[Point]bool)
	i := 0
	for line = range ReadLines(filename) {
		point.Parse(line)
		i++
		fmt.Printf("Point %v: ", point)
		for _, d := range Neighbors {
			fmt.Printf("%v ", point.Look(d))
		}
		fmt.Println("")
		if !s.Contains(point) {
			panic(fmt.Sprintf("parsing missed a point: %v", point))
		}
		correct[point] = true
	}

	for _, point = range s.boundary {
		delete(correct, point)
	}
	fmt.Printf("non-boundary points: %v\n", correct)
}

func (s *SolidShape) Attach(p Point) {
	if s.boundary == nil {
		s.boundary = make(map[Arrow]Point)
	}

	var direction Direction
	for _, direction = range Neighbors {
		var ray Arrow
		ray = p.Ray(direction)

		//fmt.Printf("Checking point=%v in direction=%v: ", p, direction)
		var end Point
		var ok bool
		end, ok = s.boundary[ray]
		if !ok || ray.Compare(p, end) == After {
			//fmt.Printf("boundary moved\n")
			s.boundary[ray] = p
			continue
		}
		//fmt.Printf("pass\n")
	}
	if !s.Contains(p) {
		panic("newly added point is not considered part of the shape")
	}
}

func (s *SolidShape) Contains(p Point) bool {
	var direction Direction
	for _, direction = range Neighbors {
		var ray Arrow
		ray = p.Ray(direction)

		var end Point
		var ok bool
		end, ok = s.boundary[ray]
		//ray.debug = p == Point{2,2,4}
		if !ok || ray.Compare(p, end) == After {
			return false
		}
	}
	return true
}

func (s *SolidShape) SurfaceArea() int {
	var total int
	seen := make(map[Point]bool)
	var p Point
	for _, p = range s.boundary {
		if seen[p] {
			continue
		}
		seen[p] = true
		total += s.surfaceAt(p)
	}
	return total
}

func (s *SolidShape) surfaceAt(p Point) int {
	fmt.Printf("Surface at %v: ", p)
	var exposed int
	var dir Direction
	var side Point
	for _, dir = range Neighbors {
		side = p.Look(dir)
		if s.Contains(side) {
			fmt.Printf("%v ", side)
			continue
		}
		exposed++
	}
	if exposed <= 0 {
		panic(fmt.Sprintf("attempting to calculate surface for an internal point %v", p))
	}
	if exposed >= 6 {
		panic(fmt.Sprintf("disjoint point %v", p))
	}
	fmt.Printf("%d\n", exposed)
	return exposed
}

func part1(filename string) string {
	shape := &SolidShape{}
	shape.FromFile(filename)
	shape.Verify(filename)
	return fmt.Sprintf("%d", shape.SurfaceArea())
}

func part2(filename string) string {
	return ""
}
