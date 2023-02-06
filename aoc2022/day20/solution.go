package main

import (
	"fmt"
)

func part1(filename string) string {
	data := ReadCoordinates(filename)
	fmt.Println(data)
	data.Mix()
	fmt.Println(data)
	return fmt.Sprint(data.Coordinates())
}

func part2(filename string) string {
	return ""
}
