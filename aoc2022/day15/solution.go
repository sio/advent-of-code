package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Point struct {
	X, Y int
}

type Rectangle struct {
	Min Point // top left corner
	Max Point // bottom right corner
}

func (r *Rectangle) Extend(p Point) {
	if p.X < r.Min.X {
		r.Min.X = p.X
	}
	if p.Y < r.Min.Y {
		r.Min.Y = p.Y
	}
	if p.X > r.Max.X {
		r.Max.X = p.X
	}
	if p.Y > r.Max.Y {
		r.Max.Y = p.Y
	}
}

func (r *Rectangle) ExtendRadius(p Point, radius int) {
	r.Extend(Point{p.X + radius, p.Y + radius})
	r.Extend(Point{p.X - radius, p.Y - radius})
}

// Manhattan distance
func (this Point) Distance(other Point) int {
	return abs(this.X-other.X) + abs(this.Y-other.Y)
}

func abs(n int) int {
	if n < 0 {
		return n * -1
	}
	return n
}

type Sensor struct {
	Location Point
	Beacon   Point
	radius   int
}

func (s *Sensor) Radius() int {
	if s.radius == 0 {
		s.radius = s.Location.Distance(s.Beacon)
	}
	return s.radius
}

func (s *Sensor) Covers(p Point) bool {
	return s.Location.Distance(p) <= s.Radius()
}

type Map struct {
	sensors  []*Sensor
	occupied map[Point]bool
	known    Rectangle
}

// Check whether a point is covered by existing sensors
// Locations of sensors and beacons themselves are not counted as covered
func (m *Map) Covered(p Point) bool {
	for _, s := range m.sensors {
		if s.Covers(p) && !m.occupied[p] {
			return true
		}
	}
	return false
}

func (m *Map) CountCovered(row int) (count int) {
	var x int
	for x = m.known.Min.X; x <= m.known.Max.X; x++ {
		if m.Covered(Point{x, row}) {
			count++
		}
	}
	return count
}

var LogFormat = regexp.MustCompile(`^Sensor at x=([0-9-]+), y=([0-9-]+): closest beacon is at x=([0-9-]+), y=([0-9-]+)$`)

func (m *Map) Parse(line string) (err error) {
	var chunks []string
	chunks = LogFormat.FindStringSubmatch(line)
	if chunks == nil || len(chunks) != 1+4 {
		return fmt.Errorf("input does not match %q regex: %s", LogFormat, line)
	}

	var numbers [4]int
	for index, chunk := range chunks[1:] {
		numbers[index], err = strconv.Atoi(chunk)
		if err != nil {
			return fmt.Errorf("could not parse a number %q from line %q: %w", chunk, line, err)
		}
	}

	sensor := &Sensor{
		Location: Point{numbers[0], numbers[1]},
		Beacon:   Point{numbers[2], numbers[3]},
	}
	m.AddSensor(sensor)
	return nil
}

func (m *Map) AddSensor(s *Sensor) {
	m.sensors = append(m.sensors, s)
	if m.occupied == nil {
		m.occupied = make(map[Point]bool)
	}
	//m.occupied[s.Location] = true
	m.occupied[s.Beacon] = true
	m.known.ExtendRadius(s.Location, s.Radius())
	m.known.Extend(s.Beacon)
}

func (m *Map) Draw() string {
	tiles := make(map[Point]rune)
	var focus *Sensor
	for _, s := range m.sensors {
		tiles[s.Location] = 'S'
		tiles[s.Beacon] = 'B'
		if s.Location == (Point{8, 7}) {
			focus = s
		}
	}
	var b strings.Builder
	var cursor Point
	for y := m.known.Min.Y; y <= m.known.Max.Y; y++ {
		for x := m.known.Min.X; x <= m.known.Max.X; x++ {
			cursor.X = x
			cursor.Y = y
			tile, found := tiles[cursor]
			switch {
			case found:
			case focus.Covers(cursor):
				tile = '#'
			default:
				tile = '.'
			}
			b.WriteRune(tile)
		}
	}
	return b.String()
}

func part1(filename string) string {
	row, found := map[string]int{
		"day15/input.txt":  2000000,
		"input.txt":        2000000,
		"day15/sample.txt": 10,
		"sample.txt":       10,
	}[filename]
	if !found {
		log.Fatalf("no target hardcoded for %s", filename)
	}
	fmt.Printf("Checking row %d\n", row)

	var err error
	cave := &Map{}
	for line := range ReadLines(filename) {
		err = cave.Parse(line)
		if err != nil {
			log.Fatalf("could not parse line: %q: %v", line, err)
		}
	}
	if filename == "sample.txt" {
		fmt.Println(cave.Draw())
	}
	return strconv.Itoa(cave.CountCovered(row))
}

func part2(filename string) string {
	return ""
}
