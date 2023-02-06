package main

import (
	"fmt"
)

func part1(filename string) string {
	monkeys := make(MonkeyGang)
	err := monkeys.Parse(filename)
	if err != nil {
		panic(err)
	}
	return fmt.Sprint(monkeys.Get("root"))
}

func part2(filename string) string {
	return ""
}
