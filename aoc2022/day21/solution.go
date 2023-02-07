package main

import (
	"fmt"
)

func part1(filename string) string {
	monkeys := MonkeyGang{}
	err := monkeys.Parse(filename)
	if err != nil {
		panic(err)
	}
	return fmt.Sprint(monkeys.Get("root"))
}

func part2(filename string) string {
	monkeys := MonkeyGang{}
	err := monkeys.Parse(filename)
	if err != nil {
		panic(err)
	}

	root := monkeys.member["root"]
	root.Job = Subtract
	human := monkeys.member["humn"]

	var iteration uint
	var oldDelta, newDelta, step MonkeyNumber
	oldDelta = monkeys.Get("root")
	step = oldDelta
	for oldDelta != 0 {
		newDelta = monkeys.Get("root")
		if newDelta == oldDelta && step != newDelta {
			step *= 10
		}
		if (newDelta > 0) != (oldDelta > 0) || abs(newDelta) > abs(oldDelta) {
			if abs(step) > 2 {
				step /= -2
			} else {
				step /= -step
			}
		}
		//fmt.Printf(
		//	"iteration=%d human=%d oldDelta=%d newDelta=%d step=%d\n",
		//	iteration,
		//	human.Number,
		//	oldDelta,
		//	newDelta,
		//	step,
		//)
		if iteration > 1000 {
			panic("could not find answer reasonably fast")
		}
		if step == 0 {
			panic("search without changing input will lead nowhere")
		}
		oldDelta = newDelta
		if oldDelta == 0 {
			break
		}
		human.Number += step
		iteration++
	}
	return fmt.Sprint(human.Number)
}

func abs(n MonkeyNumber) MonkeyNumber {
	if n < 0 {
		return n * -1
	}
	return n
}
