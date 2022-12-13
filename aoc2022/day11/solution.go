package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

type Item struct {
	Value int64
	Owner *Monkey
}

type Arithmetic int8

const (
	Add Arithmetic = 1 << iota
	Multiply
)

type InspectOperation struct {
	action Arithmetic
	arg    int64
	self   bool
}

func (op *InspectOperation) Apply(item *Item) {
	var arg int64
	arg = op.arg
	if op.self {
		arg = item.Value
	}
	switch op.action {
	default:
		panic(fmt.Sprintf("Apply() not implemented for action %b", op.action))
	case Add:
		item.Value += arg
	case Multiply:
		item.Value *= arg
	}
}

type Monkey struct {
	Business    int
	Inspection  InspectOperation
	DivideBy    int64
	Destination map[bool]int
}

type MonkeyGang struct {
	Members []*Monkey
	Items   []*Item
	Relief  bool
	Divisor int64
}

func (gang *MonkeyGang) Transfer(item *Item, owner *Monkey) {
	item.Owner = owner
}

func (gang *MonkeyGang) Obtain(item *Item) {
	gang.Items = append(gang.Items, item)
}

func (gang *MonkeyGang) Play() {
	var dest int
	var item *Item
	var monkey *Monkey
	for _, monkey = range gang.Members {
		for _, item = range gang.Items {
			if item.Owner != monkey {
				continue
			}
			monkey.Business++
			monkey.Inspection.Apply(item)
			if gang.Relief {
				item.Value /= 3 // relief
			}
			item.Value = item.Value % gang.Divisor
			dest = monkey.Destination[item.Value%monkey.DivideBy == 0]
			gang.Transfer(item, gang.Members[dest])
		}
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
		if value != len(gang.Members) {
			return fmt.Errorf("unexpected sequence number: got Monkey %d, expected %d", value, len(gang.Members))
		}
		gang.Grow()

	case strings.HasPrefix(line, PrefixItems):
		line = line[len(PrefixItems):]
		for _, chunk = range strings.Split(line, " ") {
			chunk = strings.Trim(chunk, ",")
			value, err = strconv.Atoi(chunk)
			if err != nil {
				return fmt.Errorf("could not parse a starting item: %s", chunk)
			}
			item := &Item{Value: int64(value)}
			gang.Obtain(item)
			gang.Transfer(item, gang.Last())
		}

	case strings.HasPrefix(line, PrefixOperation):
		line = line[len(PrefixOperation):]
		monkey = gang.Last()
		var chunks []string
		chunks = strings.Split(line, " ")
		if len(chunks) != 2 {
			return fmt.Errorf("unexpected inspection formula: %s (%d words)", line, len(chunks))
		}
		ops := map[string]Arithmetic{
			"+": Add,
			"*": Multiply,
		}
		action, found := ops[chunks[0]]
		if !found {
			return fmt.Errorf("unsupported arithmetic operation: %s", chunks[0])
		}
		monkey.Inspection = InspectOperation{
			action: action,
		}
		if chunks[1] == "old" {
			monkey.Inspection.self = true
		} else {
			value, err = strconv.Atoi(chunks[1])
			if err != nil {
				return fmt.Errorf("could not parse arithmetic argument: %w", err)
			}
			monkey.Inspection.arg = int64(value)
		}

	case strings.HasPrefix(line, PrefixTest):
		line = line[len(PrefixTest):]
		monkey = gang.Last()
		value, err = strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("cannot parse number: %s", line)
		}
		monkey.DivideBy = int64(value)

	case strings.HasPrefix(line, PrefixTestTrue):
		line = line[len(PrefixTestTrue):]
		monkey = gang.Last()
		value, err = strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("cannot parse number: %s", line)
		}
		monkey.Destination[true] = value

	case strings.HasPrefix(line, PrefixTestFalse):
		line = line[len(PrefixTestFalse):]
		monkey = gang.Last()
		value, err = strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("cannot parse number: %s", line)
		}
		monkey.Destination[false] = value
	}
	return nil
}

func (gang *MonkeyGang) Print() {
	fmt.Printf("Divisor: %d\n", gang.Divisor)
	for index, monkey := range gang.Members {
		fmt.Printf("Monkey %d [business=%d] holds: ", index, monkey.Business)
		for _, item := range gang.Items {
			if item.Owner != monkey {
				continue
			}
			fmt.Printf("%d ", item.Value)
		}
		fmt.Println()
	}
}

func (gang *MonkeyGang) Grow() {
	m := &Monkey{}
	m.Destination = make(map[bool]int)
	gang.Members = append(gang.Members, m)
}

func (gang *MonkeyGang) Last() *Monkey {
	return gang.Members[len(gang.Members)-1]
}

func ReadMonkeyGang(filename string) *MonkeyGang {
	gang := &MonkeyGang{}
	var err error
	for line := range ReadLines(filename) {
		err = gang.Parse(line)
		if err != nil {
			log.Fatalf("%q: %v", line, err)
		}
	}
	multipliers := make(map[int64]bool)
	for _, monkey := range gang.Members {
		multipliers[monkey.DivideBy] = true
	}
	if gang.Divisor == 0 {
		gang.Divisor = 1
	}
	for key, _ := range multipliers {
		gang.Divisor *= key
	}
	return gang
}

func (gang *MonkeyGang) PlayN(rounds int, relief bool, debug bool) int {
	gang.Relief = relief
	for i := 0; i < rounds; i++ {
		gang.Play()
	}
	if debug {
		gang.Print()
	}
	monkeys := gang.Members
	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].Business > monkeys[j].Business
	})
	return monkeys[0].Business * monkeys[1].Business
}

func part1(filename string) string {
	gang := ReadMonkeyGang(filename)
	return strconv.Itoa(gang.PlayN(20, true, false))
}

func part2(filename string) string {
	gang := ReadMonkeyGang(filename)
	return strconv.Itoa(gang.PlayN(10000, false, false))
}
