package main

import (
	"fmt"
)

func (f *Factory) QualityLevel(moves int) int {
	if f.maxGeode == 0 {
		search := searchParams{
			Limit:  moves,
			Output: ResourcePack{1, 0, 0, 0},
		}
		f.findMaxGeode(search, ResourcePack{})
	}
	fmt.Println(f)
	return f.ID * f.maxGeode
}

func (f *Factory) findMaxGeode(s searchParams, robot ResourcePack) {
	//f.Debug("[TOP %4d] building %v status %v\n", f.maxGeode, robot, s)

	// Harvest resources for this move
	s.Balance.Add(s.Output)
	if s.Balance[Geode] > f.maxGeode {
		f.maxGeode = s.Balance[Geode]
	}

	// Increase our production level thanks to new robot
	s.Output.Add(robot)

	// Check if this was our last move
	s.Limit--
	//f.Debug("           after                     %v\n", s)
	if s.Limit <= 0 {
		return
	}

	// Early exit for paths that clearly will not win
	var geodeRobot ResourcePack
	geodeRobot[Geode] = 1
	var ceiling int
	ceiling = s.Ceiling(f.Blueprint[geodeRobot])
	if s.Output[Geode] > 0 {
		f.Debug("[TOP %4d] building %v status %v\n", f.maxGeode, robot, s)
		f.Debug("           ceiling %d\n", ceiling)
		if ceiling <= f.maxGeode {
			f.Debug("           early exit\n")
		}
	}
	//f.Debug("           ceiling %d\n", ceiling)
	if ceiling <= f.maxGeode {
		return
	}

	// Launch new robot production
	var cost ResourcePack
	var ok bool
	var i int
	for i = ResourceTypeCount - 1; i >= -1; i-- { // try to build Geode robots first to increase short-circuit frequency
		robot = ResourcePack{}
		if i > 0 {
			robot[i] = 1
		}
		cost, ok = f.Blueprint[robot]
		if !ok {
			panic(fmt.Sprintf("attempting to build robot without a blueprint: %v", robot))
		}
		if !s.Balance.Affordable(cost) {
			continue
		}
		f.findMaxGeode(s.Plan(cost), robot)
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

	var newRobot, diff ResourcePack
	for s.Limit > 0 {
		s.Limit--
		newRobot = ResourcePack{}
		if s.Balance.Affordable(cost) {
			s.Balance.Spend(cost)
			max += s.Limit
			newRobot[Geode] = 1
		} else {
			diff = Sub(s.Balance, cost)
			newRobot[diff.Lowest()] = 1
		}
		s.Balance.Add(s.Output)
		s.Output.Add(newRobot)
	}
	return max
}

func part1(filename string) string {
	var factory Factory
	var result int
	for line := range ReadLines(filename) {
		factory.Parse(line)
		factory.debug = factory.ID == 2
		if !factory.debug {
			continue
		}
		fmt.Println(factory)
		result += factory.QualityLevel(24)
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
