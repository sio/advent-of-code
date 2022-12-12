package main

import (
	"fmt"
	"log"
	"math/big"
	"sort"
	"strconv"
	"strings"
)

type Item struct {
	Value *big.Int
	Owner *Monkey
}

type Arithmetic int8

const (
	Add Arithmetic = 1 << iota
	Multiply
)

type InspectOperation struct {
	action Arithmetic
	arg    *big.Int
	self   bool
}

func (op *InspectOperation) Apply(old *big.Int) *big.Int {
	var arg *big.Int
	arg = op.arg
	if op.self {
		arg = old
	}
	switch op.action {
	default:
		panic(fmt.Sprintf("Apply() not implemented for action %b", op.action))
	case Add:
		return old.Add(old, arg)
	case Multiply:
		return old.Mul(old, arg)
	}
}

type Monkey struct {
	Business        int
	operation       InspectOperation
	testDivideBy    *big.Int
	testDestination map[bool]int
}

type MonkeyGang struct {
	Members []*Monkey
	Items   []*Item
	Relief  bool
}

func (gang *MonkeyGang) Transfer(item *Item, owner *Monkey) {
	item.Owner = owner
}

func (gang *MonkeyGang) Obtain(item *Item) {
	gang.Items = append(gang.Items, item)
}

func (gang *MonkeyGang) Play() {
	var dest int
	var mod, zero *big.Int
	var item *Item
	var monkey *Monkey
	mod = big.NewInt(0)
	zero = big.NewInt(0)
	for _, monkey = range gang.Members {
		for _, item = range gang.Items {
			if item.Owner != monkey {
				continue
			}
			monkey.operation.Apply(item.Value) // inspection
			monkey.Business++
			if gang.Relief {
				item.Value.Div(item.Value, big.NewInt(3)) // relief
			}
			dest = monkey.testDestination[mod.Mod(item.Value, monkey.testDivideBy).Cmp(zero) == 0]
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
			item := &Item{Value: big.NewInt(int64(value))}
			gang.Obtain(item)
			gang.Transfer(item, gang.Last())
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
			monkey.operation.arg = big.NewInt(int64(value))
		}

	case strings.HasPrefix(line, PrefixTest):
		line = line[len(PrefixTest):]
		monkey = gang.Last()
		value, err = strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("cannot parse number: %s", line)
		}
		monkey.testDivideBy = big.NewInt(int64(value))

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
	m.testDestination = make(map[bool]int)
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
	return gang
}

func part1(filename string) string {
	gang := *ReadMonkeyGang(filename)
	gang.Relief = true
	for i := 0; i < 20; i++ {
		gang.Play()
	}
	monkeys := gang.Members
	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].Business > monkeys[j].Business
	})
	return strconv.Itoa(monkeys[0].Business * monkeys[1].Business)
}

func part2(filename string) string {
	gang := *ReadMonkeyGang(filename)
	gang.Relief = false
	for i := 0; i < 1000; i++ {
		gang.Play()
	}
	gang.Print()
	monkeys := gang.Members
	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].Business > monkeys[j].Business
	})
	return strconv.Itoa(monkeys[0].Business * monkeys[1].Business)
}
