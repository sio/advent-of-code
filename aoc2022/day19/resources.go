package main

import ()

const ResourceTypeCount = 4
const (
	Ore ResourceIndex = iota
	Clay
	Obsidian
	Geode
)
const Noop ResourceIndex = -1

var ResourceName = map[string]ResourceIndex{
	"ore":      Ore,
	"clay":     Clay,
	"obsidian": Obsidian,
	"geode":    Geode,
}

type ResourceIndex int8

type ResourcePack [ResourceTypeCount]int

func (pack *ResourcePack) Affordable(cost ResourcePack) bool {
	var i int
	for i = 0; i < len(*pack); i++ {
		if cost[i] > (*pack)[i] {
			return false
		}
	}
	return true
}

func (pack *ResourcePack) Add(income ResourcePack) {
	var i int
	for i = 0; i < len(*pack); i++ {
		(*pack)[i] += income[i]
	}
}

func (pack *ResourcePack) Spend(cost ResourcePack) {
	var i int
	for i = 0; i < len(*pack); i++ {
		(*pack)[i] -= cost[i]
	}
}

// Return index of the lowest resource
func (pack ResourcePack) Lowest() ResourceIndex {
	var minIndex ResourceIndex
	var minValue int
	minIndex = 0
	minValue = pack[minIndex]

	var i ResourceIndex
	for i = 0; i < ResourceIndex(len(pack)); i++ {
		if pack[i] < minValue {
			minValue = pack[i]
			minIndex = i
		}
	}
	return minIndex
}

// Arithmetic operations

func Sum(a, b ResourcePack) ResourcePack {
	a.Add(b)
	return a
}

func Diff(a, b ResourcePack) ResourcePack {
	a.Spend(b)
	return a
}

func Div(pack, cost ResourcePack) int {
	var nonzero bool
	var result, div int
	for i := 0; i < len(pack); i++ {
		if cost[i] == 0 {
			continue
		}
		div = pack[i] / cost[i]
		if !nonzero || div < result {
			result = div
		}
		nonzero = true
	}
	if !nonzero {
		panic("division by zero pack")
	}
	return result
}

func Mul(pack ResourcePack, multiplier int) ResourcePack {
	var i int
	for i = 0; i < len(pack); i++ {
		pack[i] *= multiplier
	}
	return pack
}
