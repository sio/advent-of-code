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

type PointIterator struct {
	Value  Point
	center Point
	radius int
	active bool
	deltaX int
	signY  int
}

func (p Point) Perimeter(radius int) *PointIterator {
	iter := &PointIterator{center: p, radius: radius}
	return iter
}

func (iter *PointIterator) Next() bool {
	if !iter.active { // loop initialization
		iter.signY = -1
		iter.deltaX = iter.radius * -1
		iter.active = true
	}
	if iter.deltaX > iter.radius && iter.signY == 1 { // loop termination
		iter.active = false
		return false
	}
	if iter.deltaX > iter.radius { // start second half circle
		iter.signY = 1
		iter.deltaX = iter.radius * -1
	}
	iter.Value.X = iter.center.X + iter.deltaX
	iter.Value.Y = iter.center.Y + iter.signY*(iter.radius-abs(iter.deltaX))
	if iter.Value.Distance(iter.center) != iter.radius {
		panic(fmt.Sprintf("point %v outside of perimeter %d for %v", iter.Value, iter.radius, iter.center))
	}
	iter.deltaX++
	return true
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
	bounds   Rectangle
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
	for x = m.bounds.Min.X; x <= m.bounds.Max.X; x++ {
		if m.Covered(Point{x, row}) {
			count++
		}
	}
	return count
}

// Beacon must be just one step out of reach of existing beacons
//
// I could not figure that out on my own, used Reddit for help (cheating)
func (m *Map) Search(min, max int) (Point, error) {
	var sensor *Sensor
	var found bool
	var iter *PointIterator
	for _, sensor = range m.sensors {
		iter = sensor.Location.Perimeter(sensor.Radius() + 1)
		for !found && iter.Next() {
			if iter.Value.X < min || iter.Value.X > max || iter.Value.Y < min || iter.Value.Y > max {
				continue
			}
			if !m.occupied[iter.Value] && !m.Covered(iter.Value) {
				found = true // we assume that only one beacon location is possible
				fmt.Printf("Found beacon: %v\n", iter.Value)
			}
		}
		if found {
			break
		}
	}
	if !found {
		return Point{}, fmt.Errorf("beacon not found in area from (%d,%d) to (%d,%d)", min, min, max, max)
	}
	return iter.Value, nil
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
	m.occupied[s.Location] = true
	m.occupied[s.Beacon] = true
	m.bounds.ExtendRadius(s.Location, s.Radius())
	m.bounds.Extend(s.Beacon)
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
	for y := m.bounds.Min.Y; y <= m.bounds.Max.Y; y++ {
		for x := m.bounds.Min.X; x <= m.bounds.Max.X; x++ {
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
		b.WriteRune('\n')
	}
	return b.String()
}

func part1(filename string) string {
	row := 2000000
	if strings.HasSuffix(filename, "sample.txt") {
		row = 10
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
	if row == 10 {
		fmt.Println(cave.Draw())
	}
	return strconv.Itoa(cave.CountCovered(row))
}

func part2(filename string) string {
	var min, max int
	min = 0
	max = 4000000
	if strings.HasSuffix(filename, "sample.txt") {
		max = 20
	}
	fmt.Printf("Checking from (%d,%d) to (%d,%d)\n", min, min, max, max)

	var err error
	cave := &Map{}
	for line := range ReadLines(filename) {
		err = cave.Parse(line)
		if err != nil {
			log.Fatalf("could not parse line: %q: %v", line, err)
		}
	}

	var beacon Point
	beacon, err = cave.Search(min, max)
	if err != nil {
		log.Fatal(err)
	}
	return strconv.Itoa(beacon.X*4000000 + beacon.Y)
}
