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
	Next      *Rope
	Trace     map[Position]bool
	DebugName string
}

func (r *Rope) Last() bool {
	return r.Next == nil
}

func (r *Rope) MoveN(delta Direction, repeat int) {
	if len(r.DebugName) != 0 {
		fmt.Printf("--- New command for %s (%d steps) ---\n", r.DebugName, repeat)
	}
	if repeat < 0 {
		panic("we don't want an endless loop!")
	}
	for i := 0; i < repeat; i++ {
		r.Move(delta)
	}
}

func (r *Rope) Move(delta Direction) {
	if len(r.DebugName) != 0 {
		fmt.Printf("move %s from %v\n", r.DebugName, r.Head)
		defer func() { fmt.Printf("move %s to   %v\n", r.DebugName, r.Head) }()
	}

	if r.Last() && r.Trace == nil { // trace only last rope in the chain
		r.Trace = make(map[Position]bool)
		r.Trace[r.Head] = true // log initial position
	}
	r.Head.Move(delta)
	if r.Last() {
		r.Trace[r.Head] = true // log new position after move
	}
	if r.Last() || r.Next.Head.Touches(r.Head) {
		return
	}

	var catchup Direction
	catchup = Direction{
		X: sign(r.Head.X - r.Next.Head.X),
		Y: sign(r.Head.Y - r.Next.Head.Y),
	}
	r.Next.Move(catchup)
	if !r.Next.Head.Touches(r.Head) {
		panic("tail did not reattach to head after move!")
	}
}

func NewRopeChain(size int) (head, tail *Rope) {
	if size <= 0 {
		panic("size must be positive")
	}

	tail = &Rope{}
	tail.Trace = make(map[Position]bool)
	tail.Trace[tail.Head] = true

	var next *Rope
	next = tail
	for i := 0; i < size; i++ {
		head = &Rope{Next: next}
		next = head
	}
	return head, tail
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

func (r *Rope) Print() {
	const size = 30
	icons := make(map[Position]rune)
	char := 'A'
	link := r
	for {
		icons[link.Head] = char
		char++
		if link.Last() {
			break
		}
		link = link.Next
	}
	fmt.Printf("\n::: Head at %v :::\n", r.Head)
	for i := -size / 2; i < size/2; i++ {
		for j := -size / 2; j < size/2; j++ {
			var found bool
			char, found = icons[Position{j, -i}]
			if !found {
				char = '.'
			}
			fmt.Printf(string(char))
		}
		fmt.Printf("\n")
	}
}

func PlaySnake(filename string, length int) string {
	motions := make(chan Motion)
	go ReadSteps(filename, motions)

	head, tail := NewRopeChain(length)
	var debug bool
	if length == 10 && strings.HasSuffix(filename, "sample2.txt") {
		debug = true
		//head.DebugName = "HEAD"
		//tail.DebugName = "TAIL"
	}
	for motion := range motions {
		if debug {
			head.Print()
		}
		head.MoveN(motion.direction, motion.repeat)
	}
	return strconv.Itoa(len(tail.Trace))
}

func part1(filename string) string {
	return PlaySnake(filename, 1)
}

func part2(filename string) string {
	return PlaySnake(filename, 10)
}
