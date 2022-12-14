package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	ListStart     = "["
	ListEnd       = "]"
	ListSeparator = ","
)

const (
	Less    int = -1
	Equal       = 0
	Greater     = 1
)

type NestedList struct {
	Value             int
	Nested            bool
	Items             []*NestedList
	Parent            *NestedList
	parsingInProgress bool
}

func IntCompare(a, b int) int {
	switch {
	case a < b:
		return Less
	case a == b:
		return Equal
	case a > b:
		return Greater
	default:
		panic("impossible branching")
	}
}

func (this *NestedList) Compare(other *NestedList) (result int) {
	//fmt.Printf("comparing %s vs %s\n", this, other)
	if !this.Nested && !other.Nested {
		//fmt.Println("ret 1")
		return IntCompare(this.Value, other.Value)
	}
	if !this.Nested && other.Nested {
		if len(other.Items) == 0 {
			//fmt.Println("ret 2")
			return Greater
		}
		result = this.Compare(other.Items[0])
		if result != Equal {
			//fmt.Println("ret 3")
			return result
		}
		if len(other.Items) > 1 {
			return Less
		}
		return Equal
	}
	if this.Nested && !other.Nested {
		//fmt.Println("ret 4")
		return other.Compare(this) * -1
	}
	var item *NestedList
	var index int
	for index, item = range this.Items {
		if index >= len(other.Items) {
			//fmt.Println("ret 5")
			return Greater
		}
		result = item.Compare(other.Items[index])
		if result != Equal {
			//fmt.Println("ret 6")
			return result
		}
	}
	switch {
	case len(this.Items) < len(other.Items):
		//fmt.Println("ret 7")
		return Less
	case len(this.Items) == len(other.Items):
		//fmt.Println("ret 8")
		return Equal
	default:
		panic("impossible branching")
	}

}

func (list *NestedList) Grow() (tail *NestedList, err error) {
	if !list.Nested {
		return nil, fmt.Errorf("cannot grow leaf node")
	}
	tail = &NestedList{Parent: list, Nested: true}
	list.Items = append(list.Items, tail)
	return tail, nil
}

func (list *NestedList) Append(value int) {
	list.Nested = true
	list.Items = append(list.Items, &NestedList{Value: value})
}

func (list *NestedList) String() string {
	if list == nil {
		return "<nil>"
	}
	if !list.Nested {
		return strconv.Itoa(list.Value)
	}

	var build strings.Builder
	var item *NestedList
	var index int
	build.WriteString(ListStart)
	for index, item = range list.Items {
		build.WriteString(item.String())
		if index < len(list.Items)-1 {
			build.WriteString(ListSeparator)
		}
	}
	build.WriteString(ListEnd)
	return build.String()
}

func (list *NestedList) Parse(line string) (err error) {
	cursor := &Cursor{list: list}
	return cursor.Parse(line)
}

type Cursor struct {
	list    *NestedList
	pointer *NestedList
}

func (cursor *Cursor) Parse(line string) (err error) {
	var chunk string
	var value int
	var closed bool
	for _, chunk = range strings.Split(line, ListSeparator) {
		switch {
		case len(chunk) == 0: // noop
		case closed:
			return fmt.Errorf("opening bracket is already closed, but received another chunk: %q", chunk)
		case strings.HasPrefix(chunk, ListStart):
			if cursor.pointer == nil {
				cursor.pointer = cursor.list
				cursor.list.Nested = true
			} else {
				cursor.pointer, err = cursor.pointer.Grow()
				if err != nil {
					return fmt.Errorf("could not grow list: %w", err)
				}
			}
			chunk = strings.TrimPrefix(chunk, ListStart)
			err = cursor.Parse(chunk)
			if err != nil {
				return fmt.Errorf("could not parse %q: %w", chunk, err)
			}
		case strings.HasSuffix(chunk, ListEnd):
			chunk = strings.TrimSuffix(chunk, ListEnd)
			err = cursor.Parse(chunk)
			if err != nil {
				return fmt.Errorf("could not parse %q: %w", chunk, err)
			}
			cursor.pointer = cursor.pointer.Parent
			if cursor.pointer == nil {
				closed = true
			}
		default:
			value, err = strconv.Atoi(chunk)
			if err != nil {
				return fmt.Errorf("invalid list value %q: %w", chunk, err)
			}
			cursor.pointer.Append(value)
		}
	}
	return nil
}

func part1(filename string) string {
	var err error
	var pair [2]*NestedList
	var index, pairIndex, result, compare int
	for line := range ReadLines(filename) {
		if len(line) == 0 {
			continue
		}
		if index > 1 || pairIndex == 0 {
			index = 0
			pairIndex++
		}
		pair[index] = &NestedList{}
		err = pair[index].Parse(line)
		if err != nil {
			return fmt.Sprintf("parsing failed: %v", err)
		}
		index++
		if index > 1 {
			compare = pair[0].Compare(pair[1])
			//fmt.Printf("#%03d: ", pairIndex)
			if compare == Less {
				result += pairIndex
				//fmt.Printf("  correct ")
			} else {
				//fmt.Printf("incorrect ")
			}
			//fmt.Printf(": %s\n                  %s (%d)\n", pair[0], pair[1], compare)
		}
	}
	return strconv.Itoa(result)
}

type Packets []*NestedList

func (p Packets) Len() int {
	return len(p)
}

func (p Packets) Less(i, j int) bool {
	return p[i].Compare(p[j]) == Less
}

func (p Packets) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Packets slice must be sorted beforehand!
func (p Packets) Find(needle ...*NestedList) (result int) {
	targets := Packets(needle)
	sort.Sort(targets)

	var targetIndex int = 0
	result = 1
	for index, packet := range p {
		if packet == targets[targetIndex] {
			targetIndex++
			result *= index + 1
			if targetIndex >= len(targets) {
				break
			}
		}
	}
	return result
}

func part2(filename string) string {
	var Needle [2]*NestedList
	for index, line := range []string{
		"[[2]]",
		"[[6]]",
	} {
		Needle[index] = &NestedList{}
		Needle[index].Parse(line)
	}

	packets := Packets(Needle[:])
	var p *NestedList
	for line := range ReadLines(filename) {
		if len(line) == 0 {
			continue
		}
		p = &NestedList{}
		p.Parse(line)
		packets = append(packets, p)
		//fmt.Printf("Got %d packets, latest: %s\n", packets.Len(), p)
	}
	sort.Sort(packets)

	var result int
	result = packets.Find(Needle[:]...)
	return strconv.Itoa(result)
}
