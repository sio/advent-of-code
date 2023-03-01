package main

import (
	"fmt"
)

func part1(filename string) string {
	maze := &Maze{}
	maze.Load(filename)
	maze.Play()
	return fmt.Sprint(maze.player.Password())
}

func part2(filename string) string {
	maze := &Maze{}
	maze.Load(filename)
	maze.ParseCube()
	maze.Play()
	return fmt.Sprint(maze.player.Password())
}
