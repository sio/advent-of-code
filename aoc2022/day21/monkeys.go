package main

import (
	"fmt"
	"strings"
)

type MonkeyNumber int64

type MonkeyGang struct {
	members map[string]Monkey
	cache   map[string]MonkeyNumber
}

func (gang *MonkeyGang) Get(name string) MonkeyNumber {
	if name == "root" { // drop cache
		gang.cache = make(map[string]MonkeyNumber)
	}

	cached, ok := gang.cache[name]
	if ok {
		return cached
	}

	monkey, ok := gang.members[name]
	if !ok {
		panic(fmt.Sprint("invalid monkey name:", name))
	}

	var left, right string
	left = monkey.Depends[0]
	right = monkey.Depends[1]

	var result MonkeyNumber
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
	gang.cache[name] = result
	return result
}

func (gang *MonkeyGang) Parse(filename string) error {
	if gang.members == nil {
		gang.members = make(map[string]Monkey)
	}

	var iter LineIterator
	if err := iter.Open(filename); err != nil {
		return err
	}
	defer iter.Close()

	for iter.Next() {
		line := strings.ReplaceAll(iter.Value(), ":", "")

		var name string
		var number MonkeyNumber

		_, err := fmt.Sscanf(line, "%s %d", &name, &number)
		if err == nil {
			gang.members[name] = Monkey{
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
		gang.members[name] = Monkey{
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
	Number  MonkeyNumber
	Depends [2]string
}

func (m Monkey) String() string {
	if m.Job == Return {
		return fmt.Sprintf("<%d>", m.Number)
	}
	return fmt.Sprintf("<%s %c %s>", m.Depends[0], m.Job, m.Depends[1])
}
