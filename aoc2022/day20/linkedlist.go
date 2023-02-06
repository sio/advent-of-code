package main

import (
	"fmt"
	"strconv"
	"strings"
)

type RingItem struct {
	Value int
	Prev  *RingItem
	Next  *RingItem
}

type Ring struct {
	Start *RingItem
	Size  int
}

func ReadCoordinates(filename string) *Ring {
	var iter LineIterator
	err := iter.Open(filename)
	if err != nil {
		panic(err)
	}
	defer iter.Close()

	var ring Ring
	var prev, first *RingItem
	for iter.Next() {
		ring.Size++
		value, err := strconv.Atoi(iter.Value())
		if err != nil {
			panic(err)
		}
		item := &RingItem{
			Value: value,
			Prev:  prev,
		}
		if prev != nil {
			prev.Next = item
		} else {
			first = item
		}
		if value == 0 {
			if ring.Start != nil {
				panic("second occurence of zero value in input")
			}
			ring.Start = item
		}
		prev = item
	}
	first.Prev = prev
	prev.Next = first
	if ring.Start == nil {
		panic("zero value not found in input")
	}
	return &ring
}

func (r *Ring) String() string {
	const showItems = 15
	var truncate string
	var show = make([]string, 0, showItems)
	item := r.Start
	for {
		show = append(show, fmt.Sprint(item.Value))
		if item.Next == r.Start {
			break
		}
		if len(show) >= showItems {
			truncate = "..."
			break
		}
		item = item.Next
	}
	return fmt.Sprintf("[Ring (%d items) %s%s]", r.Size, strings.Join(show, " "), truncate)
}
