package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

type Arithmetic int8

const (
	Add Arithmetic = 1 << iota
	Multiply
)

type InspectOperation struct {
	action Arithmetic
	arg    int
	self   bool
}

func (op *InspectOperation) Apply(old int) int {
	var arg int
	arg = op.arg
	if op.self {
		arg = old
	}
	switch op.action {
	default:
		panic(fmt.Sprintf("Apply() not implemented for action %b", op.action))
	case Add:
		return old + arg
	case Multiply:
		return old * arg
	}
}

type Monkey struct {
	Items           []int
	Business        int
	operation       InspectOperation
	testDivideBy    int
	testDestination map[bool]int
}

func (m *Monkey) Catch(item int) {
	m.Items = append(m.Items, item)
}

type MonkeyGang []*Monkey

func (gang *MonkeyGang) Play() {
	var dest int
	for _, monkey := range *gang {
		for _, item := range monkey.Items {
			item = monkey.operation.Apply(item) // inspection
			monkey.Business++
			item = item / 3 // relief
			dest = monkey.testDestination[item%monkey.testDivideBy == 0]
			(*gang)[dest].Catch(item) // throw is implied in batch when setting monkey.Items to nil
		}
		monkey.Items = nil
	}
}

const (
	PrefixMonkey    = "Monkey "
	PrefixItems     = "Starting items: "
	PrefixOperation = "Operation: new = old "
	PrefixTest      = "Test: divisible by "
	PrefixTestTrue  = "If true: throw to monkey "
	PrefixTestFalse = "If false: throw to monkey "
)

func (gang *MonkeyGang) Parse(line string) (err error) {
	line = strings.Trim(strings.TrimSpace(line), ":")
	var value int
	var chunk string
	var monkey *Monkey
	switch {

	default:
		return fmt.Errorf("parser not implemented for line: %s", line)

	case line == "":

	case strings.HasPrefix(line, PrefixMonkey):
		line = line[len(PrefixMonkey):]
		value, err = strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("could not parse monkey number: %w", err)
		}
		if value != len(*gang) {
			return fmt.Errorf("unexpected sequence number: got Monkey %d, expected %d", value, len(*gang))
		}
		gang.Grow()

	case strings.HasPrefix(line, PrefixItems):
		line = line[len(PrefixItems):]
		monkey = gang.Last()
		for _, chunk = range strings.Split(line, " ") {
			chunk = strings.Trim(chunk, ",")
			value, err = strconv.Atoi(chunk)
			if err != nil {
				return fmt.Errorf("could not parse a starting item: %s", chunk)
			}
			monkey.Catch(value)
		}

	case strings.HasPrefix(line, PrefixOperation):
		line = line[len(PrefixOperation):]
		monkey = gang.Last()
		var chunks []string
		chunks = strings.Split(line, " ")
		if len(chunks) != 2 {
			return fmt.Errorf("unexpected operation formula: %s (%d words)", line, len(chunks))
		}
		ops := map[string]Arithmetic{
			"+": Add,
			"*": Multiply,
		}
		action, found := ops[chunks[0]]
		if !found {
			return fmt.Errorf("unsupported arithmetic operation: %s", chunks[0])
		}
		monkey.operation = InspectOperation{
			action: action,
		}
		if chunks[1] == "old" {
			monkey.operation.self = true
		} else {
			value, err = strconv.Atoi(chunks[1])
			if err != nil {
				return fmt.Errorf("could not parse arithmetic argument: %w", err)
			}
			monkey.operation.arg = value
		}

	case strings.HasPrefix(line, PrefixTest):
		line = line[len(PrefixTest):]
		monkey = gang.Last()
		value, err = strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("cannot parse number: %s", line)
		}
		monkey.testDivideBy = value

	case strings.HasPrefix(line, PrefixTestTrue):
		line = line[len(PrefixTestTrue):]
		monkey = gang.Last()
		value, err = strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("cannot parse number: %s", line)
		}
		monkey.testDestination[true] = value

	case strings.HasPrefix(line, PrefixTestFalse):
		line = line[len(PrefixTestFalse):]
		monkey = gang.Last()
		value, err = strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("cannot parse number: %s", line)
		}
		monkey.testDestination[false] = value
	}
	return nil
}

func (gang *MonkeyGang) Print() {
	for index, monkey := range *gang {
		fmt.Printf("Monkey %d holds %v\n", index, monkey.Items)
	}
}

func (gang *MonkeyGang) Grow() {
	m := &Monkey{}
	m.testDestination = make(map[bool]int)
	(*gang) = append(*gang, m)
}

func (gang *MonkeyGang) Last() *Monkey {
	return (*gang)[len(*gang)-1]
}

func part1(filename string) string {
	gang := MonkeyGang{}
	var err error
	for line := range ReadLines(filename) {
		err = gang.Parse(line)
		if err != nil {
			log.Fatalf("%q: %v", line, err)
		}
	}
	for i := 0; i < 20; i++ {
		gang.Play()
	}
	sort.Slice(gang, func(i, j int) bool {
		return gang[i].Business > gang[j].Business
	})
	return strconv.Itoa(gang[0].Business * gang[1].Business)
}

func part2(filename string) string {
	return ""
}
