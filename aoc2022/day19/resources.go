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

func (this Resources) Times(n int) (result Resources) {
	result = this
	multiplier := ResourceValue(n)
	for i := 0; i < len(result); i++ {
		result[i] *= multiplier
	}
	return result
}

func (this Resources) Covers(other Resources) bool {
	for i := 0; i < len(this); i++ {
		if this[i] < other[i] {
			return false
		}
	}
	return true
}
