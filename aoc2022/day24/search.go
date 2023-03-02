package main

import (
	"fmt"
)

type SearchCursor struct {
	location Point
	round    int
}

func (cursor SearchCursor) Move(in Direction) SearchCursor {
	cursor.round++
	cursor.location = cursor.location.Look(in)
	return cursor
}

type Search struct {
	basin     *BlizzardBasin
	proximity []ScaleUnit
	cache     []PointSet
}

var Stay = Direction{0, 0}

var Moves = [...]Direction{
	Down,
	Right,
	Stay,
	Up,
	Left,
}

func (search *Search) ShortestPath() int {
	search.recurse(SearchCursor{location: search.basin.entrance})
	for index, distance := range search.proximity {
		if distance == 0 {
			return index
		}
	}
	fmt.Println("Best proximity by round:", search.proximity)
	panic("search did not reach the target, consider increasing backtrackThreshold")
}

func (search *Search) recurse(cursor SearchCursor) (ok bool) {
	// Termination condition: failure
	_, isWall := search.basin.wall[cursor.location]
	if isWall {
		return false
	}
	if search.isBlizzard(cursor) {
		return false
	}

	// Termination condition: success
	distance := cursor.location.Distance(search.basin.exit)
	if cursor.round > len(search.proximity) {
		panic("missed a proximity record in one of previous steps")
	}
	if cursor.round == len(search.proximity) {
		search.proximity = append(search.proximity, distance)
	}
	if search.proximity[cursor.round] > distance {
		search.proximity[cursor.round] = distance
	}
	if distance == 0 {
		return true
	}

	// Limit backtracking
	const backtrackThreshold = 2
	if search.proximity[cursor.round]+backtrackThreshold < distance {
		return false
	}

	// Always attempt to get closer to the target
	var next SearchCursor
	var moved bool
	for _, direction := range []Direction{Down, Right} {
		next = cursor.Move(direction)
		if search.recurse(next) {
			moved = true
		}
	}
	if moved {
		return true
	}

	// Stay in place if can not advance
	if search.recurse(cursor.Move(Stay)) {
		return true
	}

	// Backtrack if can not do anything else
	for _, direction := range []Direction{Up, Left} {
		next = cursor.Move(direction)
		if search.recurse(next) {
			moved = true
		}
	}
	return moved
}

func (search *Search) isBlizzard(cursor SearchCursor) bool {
	// Populate blizzard cache
	for !(cursor.round < len(search.cache)) {
		search.cache = append(search.cache, search.basin.Blizzards(len(search.cache)))
	}

	// Check if cursor location is blocked by a blizzard in current round
	_, blocked := search.cache[cursor.round][cursor.location]
	return blocked
}
