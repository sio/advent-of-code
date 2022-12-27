package main

import ()

const ResourceTypeCount ResourceIndex = 4
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

func (pack ResourcePack) Above(base ResourcePack) (surplus ResourcePack) {
	var i int
	for i = 0; i < len(pack); i++ {
		if pack[i] > base[i] {
			surplus[i] = pack[i] - base[i]
		}
	}
	return surplus
}

func (pack ResourcePack) Divide(piece ResourcePack) (multiplier int, ok bool) {
	var div int
	for i := 0; i < len(pack); i++ {
		if pack[i] != 0 && piece[i] == 0 {
			return 0, false // we will never reach pack in piece steps
		}
		if piece[i] == 0 {
			continue
		}
		div = pack[i] / piece[i]
		if div > multiplier {
			multiplier = div
		}
	}
	return multiplier, true
}
