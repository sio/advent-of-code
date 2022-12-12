package main

import (
	"fmt"
	"log"
	"math/big"
	"sort"
	"strconv"
	"strings"
)

var SkipRelief bool = false

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
	Items           []*big.Int
	Business        int
	operation       InspectOperation
	testDivideBy    *big.Int
	testDestination map[bool]int
}

func (m *Monkey) Catch(item *big.Int) {
	m.Items = append(m.Items, item)
}

type MonkeyGang []*Monkey

func (gang *MonkeyGang) Play() {
	var dest int
	var old, item, mod, zero *big.Int
	var monkey *Monkey
	var i int
	mod = big.NewInt(0)
	zero = big.NewInt(0)
	for i, monkey = range *gang {
		for _, old = range monkey.Items {
			item = monkey.operation.Apply(old) // inspection
			monkey.Business++
			if !SkipRelief {
				item.Div(item, big.NewInt(3)) // relief
			}
			dest = monkey.testDestination[mod.Mod(item, monkey.testDivideBy).Cmp(zero) == 0]
			(*gang)[dest].Catch(item) // throw is implied in batch when setting monkey.Items to nil
			if SkipRelief && item.Cmp(old) < 0 {
				panic(fmt.Sprintf("integer overflow from %d to %d", old, item))
			}
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
			monkey.Catch(big.NewInt(int64(value)))
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
	for index, monkey := range *gang {
		fmt.Printf("Monkey %d holds %v [business=%d]\n", index, monkey.Items, monkey.Business)
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
	gang.Print()
	for i := 0; i < 20; i++ {
		gang.Play()
		gang.Print()
	}
	sort.Slice(gang, func(i, j int) bool {
		return gang[i].Business > gang[j].Business
	})
	return strconv.Itoa(gang[0].Business * gang[1].Business)
}

func part2(filename string) string {
	return ""
	gang := *ReadMonkeyGang(filename)
	SkipRelief = true
	for i := 0; i < 1000; i++ {
		gang.Play()
	}
	//gang.Print()
	sort.Slice(gang, func(i, j int) bool {
		return gang[i].Business > gang[j].Business
	})
	return strconv.Itoa(gang[0].Business * gang[1].Business)
}
