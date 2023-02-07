package main

import (
	"fmt"
)

func part1(filename string) string {
	var maze Maze
	maze.Load(filename)
	fmt.Println(maze)
	fmt.Println(&maze)
	return ""
}

func part2(filename string) string {
	return ""
}
