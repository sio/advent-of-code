package main

import (
	"fmt"
	"sort"
)

type SearchCursor struct {
	location    Point
	destination Point
	round       int
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
	seen      map[SearchCursor]bool
}

var Stay = Direction{0, 0}

var Moves = []Direction{
	Down,
	Right,
	Stay,
	Up,
	Left,
}

func (search *Search) ShortestPath(from, to Point, startTime int) int {
	search.seen = make(map[SearchCursor]bool)

	distance := from.Distance(to)
	search.proximity = make([]ScaleUnit, startTime+1)
	for i := 0; i <= startTime; i++ {
		search.proximity[i] = distance
	}

	search.recurse(SearchCursor{
		location:    from,
		destination: to,
		round:       startTime,
	})
	for index := startTime; index < len(search.proximity); index++ {
		if search.proximity[index] == 0 {
			return index - startTime
		}
	}
	fmt.Println("Best proximity by round:", search.proximity)
	panic("search did not reach the target, consider increasing backtrackThreshold")
}

func (search *Search) recurse(cursor SearchCursor) (ok bool) {
	// Avoid infinite loops
	if search.seen[cursor] {
		return false
	}
	search.seen[cursor] = true

	// Termination condition: failure
	_, isWall := search.basin.wall[cursor.location]
	if isWall {
		return false
	}
	if search.isBlizzard(cursor) {
		return false
	}

	// Termination condition: success
	distance := cursor.location.Distance(cursor.destination)

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
	const backtrackThreshold = 30
	for i := 0; i < cursor.round-backtrackThreshold; i++ {
		if search.proximity[i] < distance {
			return false // another shorter path has achieved a significantly better result
		}
		if search.proximity[i] <= 1 {
			return false // we have already found a shorter path to the target
		}
	}

	// Try all allowed moves
	var next SearchCursor
	var moved bool
	sort.Slice(Moves, func(i, j int) bool {
		iDistance := cursor.location.Look(Moves[i]).Distance(cursor.destination)
		jDistance := cursor.location.Look(Moves[j]).Distance(cursor.destination)
		return iDistance < jDistance
	})
	for _, direction := range Moves {
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
