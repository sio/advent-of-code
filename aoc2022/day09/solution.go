package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Position struct {
	X int
	Y int
}

type Direction Position

var (
	Right Direction = Direction{1, 0}
	Left            = Direction{-1, 0}
	Up              = Direction{0, 1}
	Down            = Direction{0, -1}
)

func (this *Position) Touches(other Position) bool {
	return abs(this.X-other.X) <= 1 && abs(this.Y-other.Y) <= 1
}

func (point *Position) Move(delta Direction) {
	point.X += delta.X
	point.Y += delta.Y
}

func abs(num int) int {
	if num < 0 {
		return num * -1
	}
	return num
}

func sign(num int) int {
	switch {
	case num == 0:
		return 0
	case num > 0:
		return 1
	case num < 0:
		return -1
	default:
		panic(fmt.Sprintf("impossible branching (num=%d)", num))
	}
}

type Rope struct {
	Head      Position
	Tail      Position
	TailTrace map[Position]bool
}

func (r *Rope) MoveN(delta Direction, repeat int) {
	if repeat < 0 {
		panic("we don't want an endless loop!")
	}
	for i := 0; i < repeat; i++ {
		r.Move(delta)
	}
}

func (r *Rope) Move(delta Direction) {
	//fmt.Printf("move from %v, %v to ", r.Head, r.Tail)
	//defer func() { fmt.Printf("%v, %v\n", r.Head, r.Tail) }()

	r.Head.Move(delta)
	if r.Tail.Touches(r.Head) {
		return
	}
	var tailDelta Direction
	tailDelta = Direction{
		X: sign(r.Head.X - r.Tail.X),
		Y: sign(r.Head.Y - r.Tail.Y),
	}

	if r.TailTrace == nil {
		r.TailTrace = make(map[Position]bool)
		r.TailTrace[r.Tail] = true // log initial position
	}
	r.Tail.Move(tailDelta)
	r.TailTrace[r.Tail] = true // log new position after move
	if !r.Tail.Touches(r.Head) {
		panic("tail did not reattach to head after move!")
	}
}

type Motion struct {
	direction Direction
	repeat    int
}

func ReadSteps(filename string, motions chan<- Motion) {
	defer close(motions)
	var step Direction
	var line, command, arg string
	var repeat int
	directions := map[string]Direction{
		"R": Right,
		"L": Left,
		"U": Up,
		"D": Down,
	}
	var ok bool
	var err error
	for line = range ReadLines(filename) {
		command, arg, ok = strings.Cut(line, " ")
		if !ok {
			log.Fatalf("invalid command: %s", line)
		}
		step, ok = directions[command]
		if !ok {
			log.Fatalf("unsupported command (%s): %s", command, line)
		}
		repeat, err = strconv.Atoi(arg)
		if err != nil {
			log.Fatalf("cannot parse number of steps (%s): %s", arg, line)
		}
		motions <- Motion{direction: step, repeat: repeat}
	}
}

func part1(filename string) string {
	motions := make(chan Motion)
	go ReadSteps(filename, motions)

	rope := Rope{}
	for motion := range motions {
		//fmt.Println("--- New command ---")
		rope.MoveN(motion.direction, motion.repeat)
	}
	return strconv.Itoa(len(rope.TailTrace))
}

func part2(filename string) string {
	return ""
}
