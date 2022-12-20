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

func (vs Valves) Add(valve *Valve) Valves {
	return append(vs, valve.Name)
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
	//fmt.Printf("Calculating distance bettween %v and %v\n", a,b)
	if a.Name == b.Name {
		//fmt.Println("--short circuit")
		return 0
	}
	if a.Rate == 0 || b.Rate == 0 {
		panic(fmt.Sprintf("we assumed that zero-flow valves are not interesting!"))
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
		//fmt.Println("--known distance")
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
		//fmt.Printf("Djkstra: cursor %s\n", cursor)
		for _, valve = range cursor.Neighbors {
			//fmt.Printf("neighbor %s\n", valve)
			if visited[valve.Name] {
				continue
			}
			oldD, found = distance[valve.Name]
			newD = distance[cursor.Name] + 1 // all steps cost one minute
			if !found || newD < oldD {
				//fmt.Printf("new distance for %s: %d\n", valve, newD)
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
		if valve.Rate == 0 {
			continue // this is noise, we will never want to fetch such values
		}
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

func (g *Graph) Search(from string, depth int) int {
	start, ok := g.Get(from)
	if !ok {
		panic(fmt.Sprintf("Starting node %s not found in graph!", from))
	}
	g.DFS(start, depth, make(map[string]bool), Valves{}, 0)
	return g.MaxReward
}

func (g *Graph) DFS(start *Valve, limit int, visited map[string]bool, open Valves, collected int) {
	fmt.Printf("Top: %d | Position %s, %d remaining, open: %v, reward %d\n", g.MaxReward, start.Name, limit, open, collected)

	// termination condition
	if limit == 0 {
		if collected > g.MaxReward {
			fmt.Printf("  %04d: open valves %v\n", collected, open)
			g.MaxReward = collected
		}
		fmt.Println("  ran out of time")
		return
	}

	// early exit
	if collected+g.RewardCeiling(limit, open) <= g.MaxReward {
		fmt.Println("  this path will not win")
		return
	}

	limit--
	if start.Rate != 0 && !open.Contains(start) { // try opening this valve
		g.DFS(start, limit, visited, open.Add(start), collected+start.Rate*limit)
	}
	var valve *Valve
	for _, valve = range g.nodes { // try going to other valves instead
		g.DFS(valve, limit, visited, open, collected)
	}
	visited[start.Name] = true
}

func (g *Graph) RewardCeiling(limit int, open Valves) (max int) {
	var valve *Valve
	for _, valve = range g.nodes {
		if !open.Contains(valve) {
			max += valve.Rate * limit
		}
	}
	return max
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
	for line := range ReadLines(filename) {
		err = tunnels.Parse(line)
		if err != nil {
			log.Fatal(err)
		}
	}
	var from, to *Valve
	from, _ = tunnels.Get("CC")
	to, _ = tunnels.Get("EE")
	fmt.Println(tunnels.Distance(from, to))
	//return strconv.Itoa(tunnels.Search("AA", 30))
	return ""
}

func part2(filename string) string {
	return ""
}
