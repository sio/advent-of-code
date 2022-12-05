package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Stack []rune

func (s *Stack) Push(item rune) {
	if item == ' ' {
		return
	}
	*s = append(*s, item)
}

func (s *Stack) Pop() (item rune, ok bool) {
	if len(*s) == 0 {
		return ' ', false
	}
	index := len(*s) - 1
	item = (*s)[index]
	*s = (*s)[:index]
	return item, true
}

func (s *Stack) Top() rune {
	if len(*s) == 0 {
		return ' '
	}
	return (*s)[len(*s)-1]
}

type StackGroup []*Stack

func (sg *StackGroup) Top() string {
	tops := make([]string, len(*sg))
	for index, value := range *sg {
		tops[index] = string(value.Top())
	}
	return strings.Join(tops, "")
}

func NewStackGroup(length int) StackGroup {
	stacks := make(StackGroup, length)
	for index := range stacks {
		stacks[index] = new(Stack)
	}
	return stacks
}

type Move struct {
	Boxes int
	From  int
	To    int
}

func (m *Move) Parse(input string) {
	words := strings.Fields(input)
	if len(words) != 6 {
		log.Fatalf("incorrect command: %q", input)
	}
	command, from, to := words[0], words[2], words[4]
	if command != "move" || from != "from" || to != "to" {
		log.Fatalf("invalid command words: %q", input)
	}
	var err error
	m.Boxes, err = strconv.Atoi(words[1])
	if err != nil {
		log.Fatalf("invalid number of boxes: %q", words[1])
	}
	m.From, err = strconv.Atoi(words[3])
	if err != nil {
		log.Fatalf("invalid from address: %q", words[3])
	}
	m.To, err = strconv.Atoi(words[5])
	if err != nil {
		log.Fatalf("invalid destination address: %q", words[5])
	}
}

func (m *Move) Apply(stacks StackGroup) {
	var box rune
	var ok bool
	var source, destination *Stack
	for i := 0; i < m.Boxes; i++ {
		source = stacks[m.From-1]
		destination = stacks[m.To-1]
		box, ok = source.Pop()
		if !ok {
			log.Fatalf("could not pop a box from %v", source)
		}
		destination.Push(box)
	}
}

func part1(filename string) {
	readMoves := false
	initial := make([]string, 0)
	var stacks StackGroup
	var move Move
	for line := range ReadLines(filename) {
		if len(line) == 0 {
			continue
		}
		if readMoves {
			move.Parse(line)
			move.Apply(stacks)
			continue
		}
		if strings.HasPrefix(line, " 1   2") {
			//log.Println("Parsing initial stack configuration")
			stackLabels := strings.Fields(line)
			if len(stackLabels) > 9 {
				log.Fatalf("this implementation assumes single character stack labels")
			}
			stacks = NewStackGroup(len(stackLabels))
			for row := len(initial) - 1; row >= 0; row -= 1 {
				line = initial[row]
				//log.Printf("parsing line %q", line)
				for stackNum, stack := range stacks {
					index := stackNum*4 + 1
					runes := []rune(line)
					if index < len(runes) {
						//log.Printf("Pushing %c to stack %v", runes[index], stack)
						stack.Push(runes[index])
					}
				}
			}
			readMoves = true
			continue
		}
		initial = append(initial, line)
	}
	fmt.Printf("Part 1 result: %q", stacks.Top())
}

func part2(filename string) {}
