package main

import (
	"fmt"
	"strconv"
	"strings"
)

type RingItem struct {
	Value int64
	Prev  *RingItem
	Next  *RingItem
}

func (ri *RingItem) Look(at int64) *RingItem {
	steps := at
	if steps < 0 {
		steps = -steps
	}

	neighbor := ri
	var i int64
	for i = 0; i < steps; i++ {
		if at > 0 {
			neighbor = neighbor.Next
		} else {
			neighbor = neighbor.Prev
		}
	}
	return neighbor
}

func (ri *RingItem) Move(steps int64) {
	if steps == 0 {
		return
	}

	// Fuse links at old position
	ri.Next.Prev = ri.Prev
	ri.Prev.Next = ri.Next

	// Extend links at new position
	if steps < 0 {
		steps-- // look for the left side of the link to break
	}
	left := ri.Look(steps)
	right := left.Next
	left.Next = ri
	ri.Prev = left
	ri.Next = right
	right.Prev = ri
}

type Ring struct {
	First *RingItem
	Zero  *RingItem
	Size  int64
}

func ReadCoordinates(filename string) *Ring {
	var iter LineIterator
	err := iter.Open(filename)
	if err != nil {
		panic(err)
	}
	defer iter.Close()

	var ring Ring
	var prev *RingItem
	for iter.Next() {
		ring.Size++
		value, err := strconv.Atoi(iter.Value())
		if err != nil {
			panic(err)
		}
		item := &RingItem{
			Value: int64(value),
			Prev:  prev,
		}
		if prev != nil {
			prev.Next = item
		} else {
			ring.First = item
		}
		if value == 0 {
			if ring.Zero != nil {
				panic("second occurence of zero value in input")
			}
			ring.Zero = item
		}
		prev = item
	}
	ring.First.Prev = prev
	prev.Next = ring.First
	if ring.Zero == nil {
		panic("zero value not found in input")
	}
	return &ring
}

func (r *Ring) Append(value ...int) {
	for _, v := range value {
		item := &RingItem{Value: int64(v)}

		if r.First == nil {
			item.Next = item
			item.Prev = item
			r.First = item
			r.Size = 1
			continue
		}

		tail := r.First.Prev
		head := r.First

		item.Prev = tail
		item.Next = head

		tail.Next = item
		head.Prev = item
		r.Size++
	}
}

func (r *Ring) String() string {
	return r.StringFrom(r.First)
}

func (r *Ring) StringFrom(start *RingItem) string {
	const showItems = 15
	var truncate string
	var show = make([]string, 0, showItems)
	item := start
	for {
		show = append(show, fmt.Sprint(item.Value))
		if item.Next == start {
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

func (r *Ring) Mix() {
	r.Decrypt(1, 1)
}

func (r *Ring) Decrypt(key int64, rounds int64) {
	var cursor *RingItem

	order := make([]*RingItem, r.Size)
	cursor = r.First
	var i int64
	for i = 0; i < r.Size; i++ {
		order[i] = cursor
		cursor.Value *= key
		cursor = cursor.Next
	}

	for i = 0; i < rounds; i++ {
		for _, cursor := range order {
			steps := cursor.Value % int64(r.Size-1)
			cursor.Move(steps)
		}
	}
}

func (r *Ring) GetItem(index int64) int64 {
	index = index % r.Size
	return r.Zero.Look(index).Value
}

func (r *Ring) Coordinates() int64 {
	return r.GetItem(1000) + r.GetItem(2000) + r.GetItem(3000)
}
