package main

import (
	"fmt"
)

func part1(filename string) string {
	data := ReadCoordinates(filename)
	fmt.Println(data)
	data.Mix()
	fmt.Println(data)
	result := data.GetItem(1000) + data.GetItem(2000) + data.GetItem(3000)

	return fmt.Sprint(result)
}

func part2(filename string) string {
	return ""
}
