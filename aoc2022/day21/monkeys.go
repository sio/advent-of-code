package main

import (
	"fmt"
	"strings"
)

type MonkeyGang map[string]*Monkey

func (gang MonkeyGang) Get(name string) int64 {
	monkey, ok := gang[name]
	if !ok {
		panic(fmt.Sprint("invalid monkey name:", name))
	}

	var left, right string
	left = monkey.Depends[0]
	right = monkey.Depends[1]

	var result int64
	switch monkey.Job {
	case Return:
		return monkey.Number
	case Add:
		result = gang.Get(left) + gang.Get(right)
	case Subtract:
		result = gang.Get(left) - gang.Get(right)
	case Multiply:
		result = gang.Get(left) * gang.Get(right)
	case Divide:
		result = gang.Get(left) / gang.Get(right)
	default:
		panic(fmt.Sprintf("invalid monkey operation: %c", monkey.Job))
	}
	monkey.Job = Return
	monkey.Number = result
	return result
}

func (gang MonkeyGang) Parse(filename string) error {
	if gang == nil {
		gang = make(MonkeyGang)
	}

	var iter LineIterator
	if err := iter.Open(filename); err != nil {
		return err
	}
	defer iter.Close()

	for iter.Next() {
		line := strings.ReplaceAll(iter.Value(), ":", "")

		var name string
		var number int64

		_, err := fmt.Sscanf(line, "%s %d", &name, &number)
		if err == nil {
			gang[name] = &Monkey{
				Job:    Return,
				Number: number,
			}
			continue
		}

		var op MonkeyJob
		var left, right string
		_, err = fmt.Sscanf(line, "%s %s %c %s", &name, &left, &op, &right)
		if err != nil {
			return fmt.Errorf("could not parse line %q: %w", line, err)
		}
		gang[name] = &Monkey{
			Job:     op,
			Depends: [...]string{left, right},
		}
	}
	return nil
}

type MonkeyJob rune

const (
	Return   MonkeyJob = ' '
	Add                = '+'
	Subtract           = '-'
	Multiply           = '*'
	Divide             = '/'
)

type Monkey struct {
	Job     MonkeyJob
	Number  int64
	Depends [2]string
}

func (m Monkey) String() string {
	if m.Job == Return {
		return fmt.Sprintf("<%d>", m.Number)
	}
	return fmt.Sprintf("<%s %c %s>", m.Depends[0], m.Job, m.Depends[1])
}
