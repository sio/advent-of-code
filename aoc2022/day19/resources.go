package main

import ()

const ResourceTypeCount = 4
const (
	Ore ResourceIndex = iota
	Clay
	Obsidian
	Geode
)

var ResourceName = map[string]ResourceIndex{
	"ore":      Ore,
	"clay":     Clay,
	"obsidian": Obsidian,
	"geode":    Geode,
}

type ResourceIndex int8

type ResourcePack [ResourceTypeCount]int
