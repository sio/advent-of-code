package main

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Valves []string

func (vs *Valves) Add(valve *Valve) {
	*vs = append(*vs, valve.Name)
}

func (vs *Valves) Contains(valve *Valve) bool {
	var item string
	for _, item = range *vs {
		if valve.Name == item {
			return true
		}
	}
	return false
}

type Valve struct {
	Name      string
	Rate      int
	Neighbors []*Valve
}

func (v *Valve) String() string {
	return fmt.Sprintf("[%s: flow=%d, tunnels=%d]", v.Name, v.Rate, len(v.Neighbors))
}

var valveFormat = regexp.MustCompile(`^Valve (\w+) has flow rate=(\d+); tunnels? leads? to valves? ([\w, ]+)$`)

type Graph struct {
	nodes     map[string]*Valve
	distance  map[[2]string]int
	MaxReward int
}

func (g *Graph) Distance(a, b *Valve) int {
	if a.Name == b.Name {
		return 0
	}
	if g.distance == nil {
		g.distance = make(map[[2]string]int)
	}

	var names [2]string
	names = [...]string{a.Name, b.Name}
	sort.Strings(names[:])

	var result int
	var found bool
	result, found = g.distance[names]
	if found {
		return result
	}

	visited := make(map[string]bool)
	distance := make(map[string]int)

	var cursor, valve *Valve
	var oldD, newD int
	var name, nextName string
	var value, min int
	cursor = a
	for { // Dijkstra distances
		for _, valve = range cursor.Neighbors {
			if visited[valve.Name] {
				continue
			}
			oldD, found = distance[valve.Name]
			newD = distance[cursor.Name] + 1 // all steps cost one minute
			if !found || newD < oldD {
				distance[valve.Name] = newD
			}
		}
		visited[cursor.Name] = true

		found = false
		for name, value = range distance { // select unvisited node with lowest distance
			if visited[name] {
				continue
			}
			if !found || value < min {
				found = true
				min = value
				nextName = name
			}
		}
		if !found { // no more unvisited nodes
			break
		}
		valve, found = g.Get(nextName)
		if !found {
			panic(fmt.Sprintf("search encountered nonexistend node: %s", nextName))
		}
		if valve == cursor {
			panic(fmt.Sprintf("entering endless loop at %v", cursor))
		}
		cursor = valve
	}

	value = 0
	found = false
	for _, valve = range g.nodes { // remember all distances from node a
		names = [...]string{a.Name, valve.Name}
		sort.Strings(names[:])
		g.distance[names] = distance[valve.Name]
		if valve == b {
			value = g.distance[names]
			found = true
		}
	}
	if !found {
		panic(fmt.Sprintf("no path found between %v and %v", a, b))
	}
	return value
}

func (g *Graph) Search(from string, depth int, workers int) int {
	start, ok := g.Get(from)
	if !ok {
		panic(fmt.Sprintf("Starting node %s not found in graph!", from))
	}
	var state SearchState
	var actors SearchActors
	for i := 0; i < workers && workers > 0; i++ {
		actors.Add(SearchActor{Cursor: start, Limit: depth})
	}
	g.multiSearch(state, actors)
	return g.MaxReward
}

type SearchActor struct {
	Cursor *Valve
	Limit  int
}

type SearchActors []SearchActor

func (sa *SearchActors) Add(a SearchActor) {
	*sa = append(*sa, a)
}

func (sa *SearchActors) Cleanup(threshold int) {
	var newIndex, oldIndex int
	newIndex = 0
	for oldIndex = 0; oldIndex < len(*sa); oldIndex++ {
		if (*sa)[oldIndex].Limit < threshold {
			continue // drop this entry
		}
		(*sa)[newIndex] = (*sa)[oldIndex]
		newIndex++
	}
	(*sa) = (*sa)[:newIndex]
}

func (sa *SearchActors) Move(steps []SearchMove) (modified SearchActors) {
	modified = make(SearchActors, len(*sa))
	var i int
	for i = 0; i < len(modified); i++ {
		modified[i] = SearchActor{
			Cursor: steps[i].Dest,
			Limit:  (*sa)[i].Limit - steps[i].Cost,
		}
	}
	return modified
}

type SearchState struct {
	Path   Valves
	Reward int
}

type SearchMove struct {
	Dest *Valve
	Cost int
}

// Recursive search function for multiple actors
func (g *Graph) multiSearch(search SearchState, actors SearchActors) {

	// termination condition (can't open current valve or go anywhere)
	actors.Cleanup(1)
	if len(actors) == 0 {
		return
	}

	// early exit (this path is not a winning one)
	if search.Reward+g.RewardCeiling(search, actors) <= g.MaxReward {
		return
	}

	var i int
	for i = 0; i < len(actors); i++ {
		actors[i].Do(&search, g)
	}

	// give up early if there is nothing useful left to do
	// (one step to open next valve + at least one step to get there)
	actors.Cleanup(2)
	if len(actors) == 0 {
		return
	}

	// calculate next possible moves
	var possibilities PossibleMoves
	for i = 0; i < len(actors); i++ {
		possibilities = append(possibilities, actors[i].NextMoves(&search, g))
	}

	// drop actors without any moves left
	var keep int
	for i = 0; i < len(actors); i++ {
		if len(possibilities[i]) == 0 {
			continue
		}
		actors[keep] = actors[i]
		possibilities[keep] = possibilities[i]
		keep++
	}
	if keep == 0 {
		return
	}
	actors = actors[:keep]
	possibilities = possibilities[:keep]

	// explore all possible next moves
	var steps []SearchMove
	var permutations PossibleMovesIterator
	permutations = possibilities.Iterator()
	for permutations.Next() {
		steps = permutations.Value()
		g.multiSearch(search, actors.Move(steps))
	}
}

// Non-recursive, single iteration of search for a single actor
func (actor *SearchActor) Do(search *SearchState, g *Graph) {

	// sanity check
	if (actor.Cursor.Rate == 0 && len(search.Path) > 0) || search.Path.Contains(actor.Cursor) {
		panic(fmt.Sprintf("just did a useless move to %s", actor.Cursor))
	}

	// open current valve (except for the starting one)
	if actor.Cursor.Rate == 0 || actor.Limit < 1 {
		return
	}
	actor.Limit--
	search.Path.Add(actor.Cursor)
	search.Reward += actor.Cursor.Rate * actor.Limit

	// record results
	if search.Reward > g.MaxReward {
		g.MaxReward = search.Reward
	}

	//fmt.Printf(
	//	"Top: %4d | Position %s, %d remaining, open: %v, reward %d\n\n",
	//	g.MaxReward,
	//	actor.Cursor.Name,
	//	actor.Limit,
	//	search.Path,
	//	search.Reward,
	//)
}

// Find possible next moves
func (actor *SearchActor) NextMoves(search *SearchState, g *Graph) (moves []SearchMove) {
	var valve *Valve
	for _, valve = range g.nodes {
		if valve.Rate == 0 || search.Path.Contains(valve) || valve == actor.Cursor {
			continue
		}
		var distance int
		distance = g.Distance(actor.Cursor, valve)
		if distance+1 > actor.Limit {
			continue
		}
		moves = append(moves, SearchMove{Dest: valve, Cost: distance})
	}
	return moves
}

func (g *Graph) RewardCeiling(search SearchState, actors SearchActors) (max int) {
	var valve *Valve
	var impact, i int
	for _, valve = range g.nodes {
		if valve.Rate == 0 {
			continue
		}
		if search.Path.Contains(valve) {
			continue
		}
		impact = 0
		for i = 0; i < len(actors); i++ {
			impact = Max(impact, actors[i].Limit-g.Distance(actors[i].Cursor, valve)-1)
		}
		max += valve.Rate * impact
	}
	return max
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (g *Graph) Get(name string) (valve *Valve, ok bool) {
	valve, ok = g.nodes[name]
	return valve, ok
}

func (g *Graph) GetOrCreate(name string) *Valve {
	var ok bool
	var valve *Valve
	valve, ok = g.Get(name)
	if !ok {
		valve = &Valve{Name: name}
		if g.nodes == nil {
			g.nodes = make(map[string]*Valve)
		}
		g.nodes[valve.Name] = valve
	}
	return valve
}

func (g *Graph) ParseFile(filename string) (err error) {
	for line := range ReadLines(filename) {
		err = g.Parse(line)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Graph) Parse(line string) (err error) {
	var chunks []string
	chunks = valveFormat.FindStringSubmatch(line)
	if chunks == nil || len(chunks) != 1+3 {
		return fmt.Errorf("input does not match regex %q: %s", valveFormat, line)
	}

	var valve *Valve
	valve = g.GetOrCreate(chunks[1])
	valve.Rate, err = strconv.Atoi(chunks[2])
	if err != nil {
		return fmt.Errorf("invalid flow rate: %s (%s)", chunks[2], line)
	}
	chunks = strings.Split(chunks[3], ", ")
	var name string
	var neighbor *Valve
	for _, name = range chunks {
		neighbor = g.GetOrCreate(name)
		valve.Neighbors = append(valve.Neighbors, neighbor)
	}
	return nil
}

func part1(filename string) string {
	var err error
	tunnels := &Graph{}
	err = tunnels.ParseFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return strconv.Itoa(tunnels.Search("AA", 30, 1))
}

func part2(filename string) string {
	return ""
}
