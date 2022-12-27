package main

import (
	"fmt"
	"log"
)

const Impossible int = -99

func (f *Factory) QualityLevel(moves int) int {
	if f.maxGeode == 0 {
		branch := searchParams{
			Limit:  moves,
			Output: RobotOutput(Ore),
		}
		f.Optimization(branch, ResourcePack{})
	}
	return f.ID * f.maxGeode
}

// How many steps do we need to wait until we can build this robot?
func (f *Factory) until(robot ResourceIndex, branch searchParams) int {
	var cost, need ResourcePack
	cost = f.Blueprint[robot]

	if branch.Balance.Affordable(cost) {
		return 0
	}
	need = cost.Above(branch.Balance)
	var steps int
	var ok bool
	steps, ok = need.Divide(branch.Output)
	if !ok {
		return Impossible
	}
	return steps
}

// Shared logic for different search approaches
func (f *Factory) step(branch *searchParams, robot ResourceIndex) {
	// Harvest resources for this move
	branch.Balance.Add(branch.Output)
	if branch.Balance[Geode] > f.maxGeode {
		f.maxGeode = branch.Balance[Geode]
	}

	// Increase our production level thanks to new robot
	branch.Output.Add(RobotOutput(robot))

	// Decrease number of steps remaining
	branch.Limit--

	// Log our success
	if f.maxGeodeRobots == nil {
		f.maxGeodeRobots = make(map[int]int)
	}
	if branch.Output[Geode] > 0 {
		f.maxGeodeRobots[branch.Limit] = branch.Output[Geode]
	}
}

// Proper search for an optimal solution
func (f *Factory) Optimization(branch searchParams, robot ResourceIndex, skip int) {
	f.step(&branch, robot)
	if branch.Limit <= 0 {
		return
	}

	f.Debug("[TOP %4d] building %v status %v\n", f.maxGeode, robot, branch)
	f.Debug("%v", f.maxGeodeRobots)

	// Early exit for paths that clearly will not win
	if f.maxGeodeRobots[branch.Limit] > branch.Output[Geode] {
		f.Debug("early exit")
		return
	}
	if branch.Ceiling(f.Blueprint[Geode]) <= f.maxGeode {
		f.Debug("ceiling hit")
		return
	}

	// Launch new robot production
	var cost ResourcePack
	var ok bool
	for i := ResourceTypeCount - 1; i >= -1; i-- { // try to build Geode robots first to increase short-circuit frequency
		cost, ok = f.Blueprint[i]
		if !ok {
			panic(fmt.Sprintf("attempting to build robot without a blueprint: %v", robot))
		}
		if !branch.Balance.Affordable(cost) {
			continue
		}
		f.Optimization(branch.Plan(cost), i)
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
