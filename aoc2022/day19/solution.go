package main

import (
	"fmt"
	"log"
)

func (f *Factory) QualityLevel(moves int) int {
	if f.maxGeode == 0 {
		search := searchParams{
			Limit:  moves,
			Output: Robot(Ore),
		}
		f.OptimizationDraft(search, ResourcePack{})
		f.Optimization(search, ResourcePack{})
	}
	return f.ID * f.maxGeode
}

// Shared logic for different search approaches
func (f *Factory) step(search *searchParams, robot *ResourcePack) {
	// Harvest resources for this move
	search.Balance.Add(search.Output)
	if search.Balance[Geode] > f.maxGeode {
		f.maxGeode = search.Balance[Geode]
	}

	// Increase our production level thanks to new robot
	search.Output.Add(*robot)

	// Decrease number of steps remaining
	search.Limit--

	// Log our success
	if f.maxGeodeRobots == nil {
		f.maxGeodeRobots = make(map[int]int)
	}
	if search.Output[Geode] > 0 {
		f.maxGeodeRobots[search.Limit] = search.Output[Geode]
	}
}

// Naive first pass to populate short-circuit parameters
func (f *Factory) OptimizationDraft(search searchParams, robot ResourcePack) {
	f.step(&search, &robot)
	if search.Limit <= 0 {
		return
	}

	f.Debug("[TOP %4d] building %v status %v\n", f.maxGeode, robot, search)
	var cost, next ResourcePack
	robot = Robot(Geode)
	for {
		cost = f.Blueprint[robot]
		f.Debug("considering for next robot %v at cost %v", robot, cost)
		if search.Balance.Affordable(cost) {
			f.OptimizationDraft(search.Plan(cost), robot)
			break
		}
		next = Robot(Diff(search.Balance, cost).Lowest())
		if robot != next {
			robot = next
		} else {
			robot = Robot(Noop)
		}
	}
}

// Proper search for an optimal solution
func (f *Factory) Optimization(search searchParams, robot ResourcePack) {
	f.step(&search, &robot)
	if search.Limit <= 0 {
		return
	}

	f.Debug("[TOP %4d] building %v status %v\n", f.maxGeode, robot, search)
	f.Debug("%v", f.maxGeodeRobots)

	// Early exit for paths that clearly will not win
	if f.maxGeodeRobots[search.Limit] > search.Output[Geode] {
		f.Debug("early exit")
		return
	}
	if search.Ceiling(f.Blueprint[Robot(Geode)]) <= f.maxGeode {
		f.Debug("ceiling hit")
		return
	}

	// Launch new robot production
	var cost ResourcePack
	var ok bool
	for i := ResourceTypeCount - 1; i >= -1; i-- { // try to build Geode robots first to increase short-circuit frequency
		robot = Robot(i)
		cost, ok = f.Blueprint[robot]
		if !ok {
			panic(fmt.Sprintf("attempting to build robot without a blueprint: %v", robot))
		}
		if !search.Balance.Affordable(cost) {
			continue
		}
		f.Optimization(search.Plan(cost), robot)
		if i == Geode {
			break // we don't need to evaluate alternatives when we can afford a Geode robot
		}
	}
}

type searchParams struct {
	// How many moves we can make
	Limit int

	// Robot output for each move
	Output ResourcePack

	// What resources we have
	Balance ResourcePack
}

func (s searchParams) Plan(cost ResourcePack) searchParams {
	s.Balance.Spend(cost)
	return s
}

// Upper bound of possible Geode output for current search params
func (s searchParams) Ceiling(cost ResourcePack) int {
	var max int
	max = s.Balance[Geode] + s.Limit*s.Output[Geode]

	//var newRobot, diff ResourcePack
	for s.Limit > 0 {
		s.Limit--
		max += s.Limit
		//newRobot = ResourcePack{}
		//if s.Balance.Affordable(cost) {
		//	s.Balance.Spend(cost)
		//	max += s.Limit
		//	newRobot[Geode] = 1
		//} else {
		//	diff = Diff(s.Balance, cost)
		//	newRobot[diff.Lowest()] = 1
		//}
		//s.Balance.Add(s.Output)
		//s.Output.Add(newRobot)
	}
	return max
}

func part1(filename string) string {
	var factory Factory
	var result int
	for line := range ReadLines(filename) {
		err := factory.Parse(line)
		if err != nil {
			log.Fatal(err)
		}
		result += factory.QualityLevel(24)
		fmt.Println(factory)
	}
	return fmt.Sprintf("%d", result)
}

func part2(filename string) string {
	return ""
}

var SampleResults = [...]string{
	"33",
	"",
}
