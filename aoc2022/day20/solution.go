package main

import (
	"fmt"
)

func part1(filename string) string {
	data := ReadCoordinates(filename)
	data.Mix()
	return fmt.Sprint(data.Coordinates())
}

func part2(filename string) string {
	data := ReadCoordinates(filename)
	data.Decrypt(811589153, 10)
	return fmt.Sprint(data.Coordinates())
}
