package main

const numberOfResources = 4

type ResourceKind int

const (
	Ore ResourceKind = iota
	Clay
	Obsidian
	Geode
)

type ResourceValue int

type Resources [numberOfResources]ResourceValue

func (this Resources) Add(other Resources) (result Resources) {
	result = this
	for i := 0; i < len(result); i++ {
		result[i] += other[i]
	}
	return result
}

func (this Resources) Sub(other Resources) (result Resources) {
	result = this
	for i := 0; i < len(result); i++ {
		result[i] -= other[i]
	}
	return result
}
